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


type GatewayService interface {
    // AuthenticateToken validates the access token and returns user claims or an error.
    AuthenticateToken(token string) (map[string]interface{}, error)

    // AuthorizeRequest checks if the user has permission to access the resource.
    AuthorizeRequest(userClaims map[string]interface{}, route string, method string) error

    // RouteRequest forwards the request to the appropriate microservice based on the path and method.
    RouteRequest(path string, method string, payload []byte, headers map[string]string) ([]byte, error)

    // TransformRequest modifies the incoming request before sending it to the backend service.
    TransformRequest(path string, method string, originalPayload []byte, headers map[string]string) ([]byte, map[string]string, error)

    // TransformResponse modifies the response before returning it to the client.
    TransformResponse(responsePayload []byte, headers map[string]string) ([]byte, error)

    // LogRequest records the details of incoming requests for monitoring or debugging.
    LogRequest(userClaims map[string]interface{}, path string, method string, statusCode int)

    // HandleError provides standardized error formatting and logging for failed requests.
    HandleError(err error) ([]byte, int)
}
