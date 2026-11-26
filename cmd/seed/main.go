package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/common/hashing"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/config"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const superAdminRoleCode = "ROLE_SUPER_ADMIN"

func main() {
	viperConfig := config.NewViper()
	username := strings.TrimSpace(viperConfig.GetString("seed.admin.username"))
	if username == "" {
		username = "admin"
	}

	password := viperConfig.GetString("seed.admin.password")
	if password == "" {
		fmt.Fprintln(os.Stderr, "set seed.admin.password in config.yml before running the seeder")
		os.Exit(1)
	}

	fullName := strings.TrimSpace(viperConfig.GetString("seed.admin.full_name"))
	if fullName == "" {
		fullName = "System Administrator"
	}

	var email *string
	if value := strings.TrimSpace(viperConfig.GetString("seed.admin.email")); value != "" {
		email = &value
	}

	passwordHash, err := hashing.HashPassword(password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash admin password: %v\n", err)
		os.Exit(1)
	}

	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}
	defer connection.Close()

	err = db.Transaction(func(tx *gorm.DB) error {
		return seed(tx, username, passwordHash, fullName, email)
	})
	if err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	fmt.Printf("Seed completed. Admin username: %s\n", username)
}

func seed(tx *gorm.DB, username, passwordHash, fullName string, email *string) error {
	roles := []entity.Role{
		{Code: "ROLE_SUPER_ADMIN", Name: "Super Admin", Description: "Full access to all resources"},
		{Code: "ROLE_ADMIN", Name: "Admin", Description: "Administrative access"},
		{Code: "ROLE_USER", Name: "User", Description: "Standard user access"},
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		return fmt.Errorf("seed roles: %w", err)
	}

	permissions := dashboardPermissions()
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&permissions).Error; err != nil {
		return fmt.Errorf("seed permissions: %w", err)
	}

	var superAdmin entity.Role
	if err := tx.Where("code = ?", superAdminRoleCode).First(&superAdmin).Error; err != nil {
		return fmt.Errorf("find super admin role: %w", err)
	}

	var admin entity.User
	err := tx.Where("username = ?", username).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		admin = entity.User{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
			FullName:     fullName,
			Status:       "ACTIVE",
		}
		if err := tx.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("find admin user: %w", err)
	}

	userRole := struct {
		UserID uuid.UUID `gorm:"column:user_id"`
		RoleID uuid.UUID `gorm:"column:role_id"`
	}{UserID: admin.ID, RoleID: superAdmin.ID}
	if err := tx.Table("user_roles").Clauses(clause.OnConflict{DoNothing: true}).Create(&userRole).Error; err != nil {
		return fmt.Errorf("assign super admin role: %w", err)
	}

	if err := tx.Exec(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT ?, id FROM permissions
		ON CONFLICT DO NOTHING
	`, superAdmin.ID).Error; err != nil {
		return fmt.Errorf("assign permissions to super admin role: %w", err)
	}

	return nil
}

func dashboardPermissions() []entity.Permission {
	permissionsByResource := map[string][]string{
		"users":        {"read", "create", "update", "delete", "activate", "deactivate"},
		"merchants":    {"read", "create", "update", "delete"},
		"roles":        {"read", "create", "update", "delete"},
		"permissions":  {"read"},
		"analytics":    {"read"},
		"transactions": {"read"},
	}

	permissions := make([]entity.Permission, 0)
	for resource, actions := range permissionsByResource {
		for _, action := range actions {
			permissions = append(permissions, entity.Permission{
				Code:        resource + "." + action,
				Resource:    resource,
				Action:      action,
				Description: fmt.Sprintf("%s %s", strings.ToUpper(action), resource),
			})
		}
	}
	return permissions
}
