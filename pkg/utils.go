package pkg

import (
	"errors"
	"net/http/httputil"
	"net/url"
	"time"

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
