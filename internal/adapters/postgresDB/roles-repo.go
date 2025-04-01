package postgresDB

import (
	"fmt"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"gorm.io/gorm"
)

type RoleManagementService struct {
	dbClient *PostgresDBClient
}

// NewPostgresDBClient initializes a new GORM-based Postgres client
func NewRoleManagementRepo(dbClient *PostgresDBClient) *RoleManagementService {
	return &RoleManagementService{dbClient: dbClient}
}

// CreateRole inserts a new role into the database.
func (p *PostgresDBClient) CreateRole(role domain.Role) (domain.Role, error) {
	if err := p.DB.Create(&role).Error; err != nil {
		return domain.Role{}, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

// GetRoleByName retrieves a role by its name.
func (p *PostgresDBClient) GetRoleByName(roleName string) (domain.Role, error) {
	var role domain.Role
	err := p.DB.Where("name = ?", roleName).First(&role).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Role{}, fmt.Errorf("role '%s' not found", roleName)
		}
		return domain.Role{}, fmt.Errorf("failed to retrieve role: %w", err)
	}

	return role, nil
}

// ListRoles returns all roles from the database.
func (p *PostgresDBClient) ListRoles() ([]domain.Role, error) {
	var roles []domain.Role
	if err := p.DB.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	return roles, nil
}

// AssignRoleToUser assigns a role to a user in the user_roles table.
func (p *PostgresDBClient) AssignRoleToUser(userID string, roleName string, revocked bool) error {
	// Fetch the role by name
	dbRole, err := p.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("role '%s' not found: %w", roleName, err)
	}

	// Check if the user already has a role assigned
	var existingUserRole domain.UserRole
	err = p.DB.Where("user_id = ?", userID).First(&existingUserRole).Error

	tx := p.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err == nil {
		// Update existing role assignment
		if err := tx.Model(&existingUserRole).Update("role_id", dbRole.ID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update user role: %w", err)
		}
	} else if err == gorm.ErrRecordNotFound {
		// Assign a new role if no existing role is found
		newUserRole := domain.UserRole{UserID: userID, RoleID: dbRole.ID}
		if err := tx.Create(&newUserRole).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to assign role: %w", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("error checking existing role assignment: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
