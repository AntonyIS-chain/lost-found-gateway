package controllers

import (
	"net/http"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	service ports.AuthenticationService
}

func NewAuthenticationController(service ports.AuthenticationService) *AuthenticationController {
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
		RefreshToken string `json:"refreshToken" binding:"required"`
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
		Token string `json:"token"`
	}

	if err := ctx.ShouldBindJSON(&refreshTokenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.service.ValidateToken(refreshTokenRequest.Token)
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
		"message":    "Invalid credentials",
		"success":    false,
		"statusCode": 401,
		"valid":      res,
	})
}
