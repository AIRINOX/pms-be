package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"pms/app/http/controllers"
	"pms/app/http/middleware"
)

func Api() {
	// Authentication routes (public)
	authController := controllers.NewAuthController()
	facades.Route().Post("/auth/login", authController.Login)

	// Protected authentication routes
	facades.Route().Middleware(middleware.Auth()).Post("/auth/logout", authController.Logout)
	facades.Route().Middleware(middleware.Auth()).Get("/auth/me", authController.Me)

	// User management routes (admin only - protected by auth middleware and admin check in controller)
	userController := controllers.NewUserController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List users with pagination and search
		router.Get("/users", userController.Index)

		// Get specific user
		router.Get("/users/{id}", userController.Show)

		// Create new user
		router.Post("/users", userController.Store)

		// Update existing user
		router.Put("/users/{id}", userController.Update)

		// Delete user
		router.Delete("/users/{id}", userController.Destroy)

		// Toggle user active status
		router.Patch("/users/{id}/toggle-status", userController.ToggleStatus)
	})

	// Role management routes (admin only)
	roleController := controllers.NewRoleController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List all roles
		router.Get("/roles", roleController.Index)
	})
}
