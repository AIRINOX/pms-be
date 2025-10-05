package seeders

import (
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"

	"pms/app/models"
)

type UserSeeder struct {
}

// Signature The unique signature for the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder.
func (s *UserSeeder) Run() error {
	// Check if admin user already exists
	var existingUser models.User
	if err := facades.Orm().Query().Where("username", "admin").FirstOrFail(&existingUser); err == nil {
		facades.Log().Info(err)
		facades.Log().Info(existingUser.ID)
		facades.Log().Info("Admin user already exists, skipping seeder")
		return nil
	}

	// Get admin role
	var adminRole models.Role
	if err := facades.Orm().Query().Where("key", "admin").FirstOrFail(&adminRole); err != nil {
		facades.Log().Error("Admin role not found. Please run RoleSeeder first.")
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		facades.Log().Error("Failed to hash admin password")
		return err
	}

	// Create admin user
	adminUser := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Name:     "System Administrator",
		Email:    "admin@pms-airinox.com",
		Phone:    "",
		RoleID:   adminRole.ID,
		IsActive: true,
	}

	if err := facades.Orm().Query().Create(&adminUser); err != nil {
		facades.Log().Error("Failed to create admin user")
		return err
	}

	facades.Log().Info("Created admin user with username: admin, password: admin123")
	return nil
}
