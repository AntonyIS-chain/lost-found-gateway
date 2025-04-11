package postgresDB

import (
	"errors"
	"fmt"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationManagementRepo struct {
	db *gorm.DB
}

// NewAuthenticationManagementRepo initializes the repository with a GORM DB client
func NewAuthenticationManagementRepo(db *gorm.DB) *AuthenticationManagementRepo {
	return &AuthenticationManagementRepo{db: db}
}

func (c *PostgresDBClient) RegisterUser(user domain.User) (domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	if err := c.DB.Create(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("failed to register user: %w", err)
	}

	userToken := domain.UserToken{ID: user.ID, Token: ""}
	
	if err := c.DB.Create(&userToken).Error; err != nil {
		return domain.User{}, fmt.Errorf("failed to create user id and token: %w", err)
	}

	return user, nil
}

// GetUserByEmail fetches a user by email and verifies the password
func (c *PostgresDBClient) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User
	if err := c.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return domain.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (c *PostgresDBClient) IsRefreshTokenValid(ID, refreshToken string) (bool, error) {
	var token domain.UserToken
	err := c.DB.Select("id").Where("id = ? AND token = ? AND revoked = false", ID, refreshToken).Take(&token).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	return true, nil
}

func (c *PostgresDBClient) StoreRefreshToken(ID, refreshToken string) error {
	userToken := domain.UserToken{}
	err := c.DB.Model(&userToken).Where("id = ?", ID).Updates(map[string]interface{}{
		"token":   refreshToken,
		"revoked": false,
	}).Error

	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

func (c *PostgresDBClient) InvalidateRefreshToken(refreshToken string) error {

	userID, err := pkg.ExtractUserIDFromToken(refreshToken)

	if err != nil {
		return fmt.Errorf("failed to extract userID from token: %w", err)
	}
	// Set the token as revoked instead of deleting it (better for logging)
	result := c.DB.Model(&domain.UserToken{}).
		Where("id = ? AND token = ?", userID, refreshToken).
		Update("revoked", true) // Assuming you have a `revoked` column
	if result.Error != nil {
		return fmt.Errorf("failed to invalidate refresh token: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}

	return nil
}

// ListUsers retrieves all users with their roles
func (c *PostgresDBClient) ListUsers() ([]domain.User, error) {
	var users []domain.User

	if err := c.DB.Preload("Role").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
