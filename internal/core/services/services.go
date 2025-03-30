package services

import (
	"errors"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"golang.org/x/crypto/bcrypt"
)

// JWT Secret Keys (use environment variables in production)
var accessTokenSecret = []byte("your-access-secret-key")
var refreshTokenSecret = []byte("your-refresh-secret-key")

type AuthenticationManagementService struct {
	repo   ports.AuthenticationRepository
	config config.Config
}

func NewAuthenticationManagementService(repo ports.AuthenticationRepository, config config.Config) *AuthenticationManagementService {
	return &AuthenticationManagementService{
		repo:   repo,
		config: config,
	}
}

func (a *AuthenticationManagementService) Authenticate(email, password string) (domain.AuthenticationResponse, error) {
	user, err := a.repo.GetUserByEmail(email)

	if err != nil {
		return domain.AuthenticationResponse{}, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.AuthenticationResponse{}, errors.New("invalid email or password")
	}

	accessToken, err := pkg.GenerateToken(user.ID, []byte(a.config.SECRET_KEY), 15*time.Minute)
	if err != nil {
		return domain.AuthenticationResponse{}, err
	}

	refreshToken, err := pkg.GenerateToken(user.ID, refreshTokenSecret, 7*24*time.Hour)
	if err != nil {
		return domain.AuthenticationResponse{}, err
	}

	err = a.repo.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		return domain.AuthenticationResponse{}, errors.New("failed to store refresh token")
	}

	return domain.AuthenticationResponse{Access_token: accessToken, Refresh_token: refreshToken}, nil
}

// RefreshToken generates a new access token if the refresh token is valid
func (a *AuthenticationManagementService) RefreshToken(refreshToken string) (domain.AuthenticationResponse, error) {
	// Validate the refresh token
	claims, err := pkg.ValidateToken(refreshToken, refreshTokenSecret)
	if err != nil {
		return domain.AuthenticationResponse{}, errors.New("invalid refresh token")
	}

	ID, ok := claims["ID"].(string)

	if !ok {
		return domain.AuthenticationResponse{}, errors.New("invalid token data")
	}

	// Check if the refresh token exists in DB
	valid, err := a.repo.IsRefreshTokenValid(ID, refreshToken)
	if err != nil || !valid {
		return domain.AuthenticationResponse{}, errors.New("refresh token revoked or not found")
	}

	// Generate new access token
	accessToken, err := pkg.GenerateToken(ID, []byte(a.config.SECRET_KEY), 15*time.Minute)
	if err != nil {
		return domain.AuthenticationResponse{}, err
	}

	// Generate new refresh token
	newRefreshToken, err := pkg.GenerateToken(ID, refreshTokenSecret, 7*24*time.Hour)
	if err != nil {
		return domain.AuthenticationResponse{}, err
	}

	// Store new refresh token in DB
	err = a.repo.StoreRefreshToken(ID, newRefreshToken)
	if err != nil {
		return domain.AuthenticationResponse{}, errors.New("failed to store new refresh token")
	}

	return domain.AuthenticationResponse{Access_token: accessToken, Refresh_token: newRefreshToken}, nil
}

// ValidateToken checks if the given token is valid and not expired
func (a *AuthenticationManagementService) ValidateToken(token string) (bool, error) {
	_, err := pkg.ValidateToken(token, accessTokenSecret)
	if err != nil {
		return false, errors.New("invalid token")
	}
	return true, nil
}
