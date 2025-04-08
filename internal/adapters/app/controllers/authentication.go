package controllers

import (
	"net/http"
	"strings"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	service ports.AuthService
}

func NewAuthenticationController(service ports.AuthService) *AuthenticationController {
	return &AuthenticationController{service: service}
}

func (uc *AuthenticationController) Authenticate(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "success": false})
		return
	}

	response, err := uc.service.Authenticate(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "success": false})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (uc *AuthenticationController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "success": false})
		return
	}

	response, err := uc.service.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token", "success": false})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (uc *AuthenticationController) ValidateToken(ctx *gin.Context) {
	var refreshTokenRequest struct {
		Token string `json:"access_token"`
	}

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if strings.TrimSpace(refreshTokenRequest.Token) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "access_token is required"})
		return
	}

	claims, err := uc.service.ValidateToken(refreshTokenRequest.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message":    err.Error(),
			"success":    false,
			"statusCode": http.StatusUnauthorized,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Valid credentials",
		"success":    true,
		"statusCode": http.StatusOK,
		"valid":      true,
		"claims":     claims,
	})
}

func (uc *AuthenticationController) InValidateToken(ctx *gin.Context) {
	var refreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.service.InvalidateRefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message":    "Invalid credentials",
			"success":    false,
			"statusCode": 401,
		})
		return
	}

	// Return tokens
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Token revocked successfuly",
		"success":    true,
		"statusCode": 200,
	})
}
