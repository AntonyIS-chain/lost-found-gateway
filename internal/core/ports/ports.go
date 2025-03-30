package ports

import "github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"

type AuthenticationService interface {
	Authenticate(email,password string) (response domain.AuthenticationResponse, err error)
	RefreshToken(refreshToken string) (response domain.AuthenticationResponse, err error)
	ValidateToken(token string) (bool, error)
}

type AuthenticationRepository interface {
	GetUserByEmail(email string) (user domain.User, err error)
	StoreRefreshToken(userID, refreshToken string) error
	IsRefreshTokenValid(userID, refreshToken string) (bool, error)
}

type ProxyService interface {
	ForwardRequest(req domain.ProxyRequest) (domain.ProxyResponse, error)
}

type AuthService interface {
	Authenticate(username, password string) (string, error)
	ValidateToken(token string) (bool, error)
}
