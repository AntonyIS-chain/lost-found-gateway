package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func (a *AuthenticationManagementService) Authenticate(email, password string) (domain.AuthResponse, error) {
	user, err := a.repo.GetUserByEmail(email)

	if err != nil {
		return domain.AuthResponse{}, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.AuthResponse{}, errors.New("invalid email or password")
	}

	accessToken, err := pkg.GenerateToken(user.ID, []byte(a.config.ACCESS_TOKEN_SECRET_KEY), 15*time.Minute)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	refreshToken, err := pkg.GenerateToken(user.ID, []byte(a.config.REFRESH_TOKEN_SECRET_KEY), 7*24*time.Hour)
	if err != nil {
		return domain.AuthResponse{}, err
	}

	err = a.repo.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		return domain.AuthResponse{}, errors.New("failed to store refresh token")
	}

	return domain.AuthResponse{Access_token: accessToken, Refresh_token: refreshToken}, nil
}

// RefreshToken generates a new access token if the refresh token is valid
func (a *AuthenticationManagementService) RefreshToken(refreshToken string) (domain.AuthResponse, error) {
	// Validate the refresh token using the correct secret key
	claims, err := pkg.ValidateToken(refreshToken, []byte(a.config.REFRESH_TOKEN_SECRET_KEY))
	if err != nil {
		return domain.AuthResponse{}, errors.New("invalid refresh token")
	}

	// Extract the user ID from claims (check if the key should be "id" or "sub")
	ID, ok := claims["id"].(string)
	if !ok {
		fmt.Println("Invalid token data, ID missing:", claims)
		return domain.AuthResponse{}, errors.New("invalid token data")
	}

	// Check if the refresh token exists in DB
	valid, err := a.repo.IsRefreshTokenValid(ID, refreshToken)
	if err != nil || !valid {
		fmt.Println("Refresh token revoked or not found in DB")
		return domain.AuthResponse{}, errors.New("refresh token revoked or not found")
	}

	// Generate new access token
	accessToken, err := pkg.GenerateToken(ID, []byte(a.config.ACCESS_TOKEN_SECRET_KEY), 15*time.Minute)
	if err != nil {
		fmt.Println("Failed to generate access token:", err)
		return domain.AuthResponse{}, err
	}

	// Generate new refresh token
	newRefreshToken, err := pkg.GenerateToken(ID, []byte(a.config.REFRESH_TOKEN_SECRET_KEY), 7*24*time.Hour)
	if err != nil {
		fmt.Println("Failed to generate new refresh token:", err)
		return domain.AuthResponse{}, err
	}

	// Store the new refresh token in the database
	err = a.repo.StoreRefreshToken(ID, newRefreshToken)
	if err != nil {
		fmt.Println("Failed to store new refresh token:", err)
		return domain.AuthResponse{}, errors.New("failed to store new refresh token")
	}

	// Return the new tokens
	return domain.AuthResponse{
		Access_token:  accessToken,
		Refresh_token: newRefreshToken,
	}, nil
}

// ValidateToken checks if the given token is valid and not expired
func (a *AuthenticationManagementService) ValidateToken(token string) (bool, error) {
	_, err := pkg.ValidateToken(token, []byte(a.config.SECRET_KEY))
	if err != nil {
		return false, errors.New("invalid token")
	}
	return true, nil
}

// InvalidateRefreshToken revokes the given refresh token (logout)
func (a *AuthenticationManagementService) InvalidateRefreshToken(refreshToken string) error {
	
	// Call the repository method to revoke the token
	err := a.repo.InvalidateRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to invalidate refresh token: %w", err)
	}

	return nil
}

func (s *AuthenticationManagementService) ListUsers() ([]domain.User, error) {
	return s.repo.ListUsers()
}

func (s *AuthenticationManagementService) RegisterUser(user domain.User) (domain.User, error) {
	// Check if user already exists by email
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err == nil && existingUser.ID != "" {
		return domain.User{}, err
	}

	// Assign RoleID instead of Role struct
	userId := uuid.New().String()
	user.ID = userId

	return s.repo.RegisterUser(user)
}
