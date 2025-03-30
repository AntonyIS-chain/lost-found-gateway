package postgresDB

import (
	"fmt"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"gorm.io/gorm"
)

type AuthenticationManagementRepo struct {
	db *gorm.DB
}

// NewAuthenticationManagementRepo initializes the repository with a GORM DB client
func NewAuthenticationManagementRepo(db *gorm.DB) *AuthenticationManagementRepo {
	return &AuthenticationManagementRepo{db: db}
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
	var count int64
	err := c.DB.Model(&domain.UserToken{}).Where("id = ? AND token = ?", ID, refreshToken).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (c *PostgresDBClient) StoreRefreshToken(ID, refreshToken string) error {
	userToken := domain.UserToken{}
	err := c.DB.Model(&userToken).Where("id = ?", ID).Update("token", refreshToken).Error
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}
