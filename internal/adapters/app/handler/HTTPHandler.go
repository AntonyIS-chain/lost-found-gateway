package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
)

type GatewayHandler struct {
	service ports.GatewayService
	conf    *config.Config
}

// NewGatewayHandler initializes a new GatewayHandler
func NewGatewayHandler(service ports.GatewayService) *GatewayHandler {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return &GatewayHandler{
		service: service,
		conf:    conf,
	}
}

// ServeHTTP intercepts all incoming HTTP requests and routes or proxies them accordingly
func (h *GatewayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("➡️  Incoming Request: %s %s", r.Method, r.URL.Path)

	if isPublicRoute(r.URL.Path) {
		h.proxyToService(w, r)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg, code := h.service.HandleError(http.ErrNoCookie)
		http.Error(w, string(msg), code)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.service.AuthenticateToken(token)
	if err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	if err := h.service.AuthorizeRequest(claims, r.URL.Path, r.Method); err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	payload, headers, err := h.service.TransformRequest(r.URL.Path, r.Method, body, getHeadersMap(r.Header))
	if err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	response, err := h.service.RouteRequest(r.URL.Path, r.Method, payload, headers)
	if err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	finalResp, err := h.service.TransformResponse(response, headers)
	if err != nil {
		msg, code := h.service.HandleError(err)
		http.Error(w, string(msg), code)
		return
	}

	h.service.LogRequest(claims, r.URL.Path, r.Method, http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalResp)
}

// Public routes handled here
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/v1/auth/signin",
		"/v1/auth/signup",
		"/v1/auth/forgot",
		"/v1/auth/reset",
		"/v1/matching",
	}
	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// Map URL prefix to corresponding microservice URL from config
func (h *GatewayHandler) getServiceURL(path string) string {
	switch {
	case strings.HasPrefix(path, "/v1/auth"):
		return h.conf.USER_SERVICE
	case strings.HasPrefix(path, "/v1/rewards"):
		return h.conf.REWARD_SERVICE
	case strings.HasPrefix(path, "/v1/payments"):
		return h.conf.PAYMENT_SERVICE
	case strings.HasPrefix(path, "/v1/notifications"):
		return h.conf.NOTIFICATION_SERVICE
	case strings.HasPrefix(path, "/v1/matching"):
		return h.conf.MATCHING_SERVICE
	case strings.HasPrefix(path, "/v1/documents"):
		return h.conf.DOCUMENT_SERVICE
	default:
		return "" // Unknown service
	}
}

// Generic proxy method
func (h *GatewayHandler) proxyToService(w http.ResponseWriter, r *http.Request) {
	targetBase := h.getServiceURL(r.URL.Path)
	if targetBase == "" {
		http.Error(w, "No matching service found", http.StatusNotFound)
		return
	}

	targetURL := fmt.Sprintf("%s%s?%s", targetBase, r.URL.Path, r.URL.RawQuery)
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	req.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Target service unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// Converts headers to flat map
func getHeadersMap(h http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range h {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
