package domain

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID           string `json:"id" db:"id"`
	FirstName    string `json:"first_name" db:"first_name"`
	LastName     string `json:"last_name" db:"last_name"`
	Email        string `json:"email"  db:"email"`
	Token        string `json:"token"  db:"token"`
	Phone        string `json:"phone,omitempty" db:"phone"`
	PasswordHash string `json:"-" db:"password_hash"`
}

type UserToken struct {
	ID           string `json:"id" db:"id"`
	Token        string `json:"token"  db:"token"`
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

type AuthenticationResponse struct {
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}
