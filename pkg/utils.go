package pkg

import (
	"errors"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// NewReverseProxy creates a reverse proxy to forward requests
func NewReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func GenerateToken(id string, secret []byte, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration)
	claims := jwt.MapClaims{
		"id":  id,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// Helper function to validate JWT token
func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ExtractUserIDFromToken(refreshToken string) (string, error) {
	
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		conf, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		// Provide the secret key used to sign the token
		return []byte(conf.REFRESH_TOKEN_SECRET_KEY), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if id, exists := claims["id"].(string); exists {
			return id, nil
		}
		return "", fmt.Errorf("user_id not found in token")
	}

	return "", fmt.Errorf("invalid token claims")
}
