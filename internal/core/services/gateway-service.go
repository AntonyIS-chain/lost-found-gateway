package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
)


type gatewayServiceImpl struct{}

func NewGatewayService() ports.GatewayService {
	return &gatewayServiceImpl{}
}

func (g *gatewayServiceImpl) AuthenticateToken(tokenStr string) (map[string]interface{}, error) {
	// Remove Bearer prefix if present
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}


	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		conf, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}


		return conf.SECRET_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert claims to map[string]interface{}
		userClaims := make(map[string]interface{})
		for k, v := range claims {
			userClaims[k] = v
		}
		return userClaims, nil
	}

	return nil, errors.New("invalid token")
}

func (g *gatewayServiceImpl) AuthorizeRequest(userClaims map[string]interface{}, route string, method string) error {
	role, ok := userClaims["role"].(string)
	if !ok || role == "" {
		return errors.New("missing role in token claims")
	}

	// Example rule: only admin can DELETE
	if method == http.MethodDelete && role != "admin" {
		return errors.New("unauthorized: insufficient permissions")
	}

	return nil
}

func (g *gatewayServiceImpl) RouteRequest(path string, method string, payload []byte, headers map[string]string) ([]byte, error) {
	// Stub: In a real gateway, you'd proxy the request using http.Client or a load balancer
	log.Printf("Routing request to path: %s method: %s", path, method)
	return []byte(fmt.Sprintf(`{"status":"success","path":"%s"}`, path)), nil
}

func (g *gatewayServiceImpl) TransformRequest(path string, method string, originalPayload []byte, headers map[string]string) ([]byte, map[string]string, error) {
	// Example: inject correlation ID
	headers["X-Correlation-ID"] = fmt.Sprintf("req-%d", time.Now().UnixNano())
	return originalPayload, headers, nil
}

func (g *gatewayServiceImpl) TransformResponse(responsePayload []byte, headers map[string]string) ([]byte, error) {
	// Example: add a response footer or metadata
	return append(responsePayload, []byte(`,"gateway":"ok"}`)...), nil
}

func (g *gatewayServiceImpl) LogRequest(userClaims map[string]interface{}, path string, method string, statusCode int) {
	log.Printf("User: %v | Path: %s | Method: %s | Status: %d", userClaims["sub"], path, method, statusCode)
}

func (g *gatewayServiceImpl) HandleError(err error) ([]byte, int) {
	log.Printf("Gateway error: %v", err)
	return []byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())), http.StatusUnauthorized
}
