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

	// Client management routes (commercial/admin only)
	clientController := controllers.NewClientController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List clients with pagination, search and filtering
		router.Get("/clients", clientController.Index)

		// Get specific client with sites
		router.Get("/clients/{clientId}", clientController.Show)

		// Create new client
		router.Post("/clients", clientController.Store)

		// Update existing client
		router.Put("/clients/{clientId}", clientController.Update)

		// Delete client
		router.Delete("/clients/{clientId}", clientController.Destroy)

		// Get all orders for a specific client
		router.Get("/clients/{clientId}/orders", clientController.GetClientOrders)
	})

	// Client site management routes (commercial/admin only)
	clientSiteController := controllers.NewClientSiteController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List sites for a specific client
		router.Get("/clients/{clientId}/sites", clientSiteController.Index)

		// Get specific client site
		router.Get("/clients/{clientId}/sites/{siteId}", clientSiteController.Show)

		// Create new client site
		router.Post("/clients/{clientId}/sites", clientSiteController.Store)

		// Update existing client site
		router.Put("/clients/{clientId}/sites/{siteId}", clientSiteController.Update)

		// Delete client site
		router.Delete("/clients/{clientId}/sites/{siteId}", clientSiteController.Destroy)
	})
}
