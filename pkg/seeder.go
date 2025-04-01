package pkg

import (
	"log"

	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/ports"
)

func SeedRoles(roleService ports.RoleService) {
	existingRoles, err := roleService.ListRoles()
	if err != nil {
		log.Printf("Failed to fetch roles: %v", err)
		return
	}

	roleNames := []domain.Role{
		{Name: "User Admin", Description: "Administrator with full access to the system"},
		{Name: "Moderator", Description: "Can moderate user content"},
		{Name: "Registered User", Description: "Standard user with basic privileges"},
		{Name: "Guest User", Description: "Limited access user"},
	}

	existingRoleMap := make(map[string]bool)
	for _, role := range existingRoles {
		existingRoleMap[role.Name] = true
	}

	for _, role := range roleNames {
		if existingRoleMap[role.Name] {
			// log.Printf("Role '%s' already exists. Skipping seeding.\n", role.Name)
			continue
		}

		_, err := roleService.CreateRole(role)
		if err != nil {
			log.Printf("Failed to create role '%s': %v", role.Name, err)
			continue
		}
	}

	log.Println("Successfully seeded all roles.")
}

// SeedUsers populates the database with initial users and roles
func SeedUsers(userService ports.AuthService, roleService ports.RoleService) {
	existingUsers, err := userService.ListUsers()
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return
	}

	existingUserMap := make(map[string]bool)
	for _, user := range existingUsers {
		existingUserMap[user.Email] = true
	}

	existingRoles, err := roleService.ListRoles()
	if err != nil {
		log.Printf("Failed to fetch roles: %v", err)
		return
	}

	roleMap := make(map[string]int)
	for _, role := range existingRoles {
		roleMap[role.Name] = role.ID
	}

	users := []domain.User{
		{FirstName: "Admin", LastName: "User", Email: "admin@example.com", RoleName: "User Admin"},
		{FirstName: "Moderator", LastName: "User", Email: "moderator@example.com", RoleName: "Moderator"},
		{FirstName: "John", LastName: "Doe", Email: "user@example.com", RoleName: "Registered User"},
		{FirstName: "Guest", LastName: "User", Email: "guest@example.com", RoleName: "Guest User"},
	}

	for _, user := range users {
		if existingUserMap[user.Email] {
			continue
		}

		// Validate role existence
		roleID, exists := roleMap[user.RoleName]
		if !exists {
			log.Printf("Role '%s' not found. Skipping user %s.\n", user.RoleName, user.Email)
			continue
		}
		user.RoleID = roleID
		user.PasswordHash = "Password@1234"

		createdUser, err := userService.RegisterUser(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
			continue
		}

		// Assign role after user creation
		err = roleService.AssignRoleToUser(createdUser.ID, user.RoleName, false)
		if err != nil {
			log.Printf("Failed to assign role to user %s: %v", user.Email, err)
			continue
		}
	}

	log.Println("Successfully seeded all users and assigned roles.")
}
