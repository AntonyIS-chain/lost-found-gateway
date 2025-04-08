package domain

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID           string    `json:"id" db:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Email        string    `json:"email"  db:"email"`
	Phone        string    `json:"phone,omitempty" db:"phone"`
	PasswordHash string    `json:"-" db:"password_hash"`
	RoleID       int       `json:"role_id" db:"role_id"`
	Role         Role      `json:"role" db:"-"`
	RoleName     string    `json:"role_name" db:"role_name"`
	IsActive     bool      `json:"is_active" db:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserToken struct {
	ID      string `json:"id" gorm:"column:id"`
	Token   string `json:"token" gorm:"column:token"`
	Revoked bool   `json:"revoked" gorm:"column:revoked"`
}


type Role struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UserRole struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"not null;index"`
	RoleID int    `gorm:"not null;index"`
}

type ProxyRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type ProxyResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type AuthResponse struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}
