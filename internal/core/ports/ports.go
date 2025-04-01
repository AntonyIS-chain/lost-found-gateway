package ports

import "github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"

type AuthService interface {
	Authenticate(email, password string) (response domain.AuthResponse, err error)
	RefreshToken(refreshToken string) (response domain.AuthResponse, err error)
	ValidateToken(token string) (bool, error)
	InvalidateRefreshToken(token string) error
	ListUsers() ([]domain.User, error)
	RegisterUser(user domain.User) (domain.User, error)
}

type AuthenticationRepository interface {
	GetUserByEmail(email string) (user domain.User, err error)
	StoreRefreshToken(userID, refreshToken string) error
	IsRefreshTokenValid(userID, refreshToken string) (bool, error)
	InvalidateRefreshToken(token string) error
	ListUsers() ([]domain.User, error)
	RegisterUser(user domain.User) (domain.User, error)
}

type RoleRepository interface {
	CreateRole(role domain.Role) (domain.Role, error)
	GetRoleByName(roleName string) (domain.Role, error)
	ListRoles() ([]domain.Role, error)
	AssignRoleToUser(userID string, roleName string, revocked bool) error
}

type RoleService interface {
	CreateRole(role domain.Role) (domain.Role, error)
	GetRoleByName(roleName string) (domain.Role, error)
	ListRoles() ([]domain.Role, error)
	AssignRoleToUser(userID string, roleName string, revocked bool) error
}

type ProxyService interface {
	ForwardRequest(req domain.ProxyRequest) (domain.ProxyResponse, error)
}
