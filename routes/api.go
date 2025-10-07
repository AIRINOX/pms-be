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

	// Storage Location management routes (methodes/admin only)
	storageLocationController := controllers.NewStorageLocationController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List storage locations with pagination, search and filtering
		router.Get("/storage-locations", storageLocationController.Index)

		// Get specific storage location
		router.Get("/storage-locations/{id}", storageLocationController.Show)

		// Create new storage location
		router.Post("/storage-locations", storageLocationController.Store)

		// Update existing storage location
		router.Put("/storage-locations/{id}", storageLocationController.Update)

		// Delete storage location
		router.Delete("/storage-locations/{id}", storageLocationController.Destroy)
	})

	// Category management routes (methodes/admin only)
	categoryController := controllers.NewCategoryController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List categories with pagination, search and filtering
		router.Get("/categories", categoryController.Index)

		// Get categories tree structure
		router.Get("/categories/tree", categoryController.GetTree)

		// Get specific category
		router.Get("/categories/{id}", categoryController.Show)

		// Create new category
		router.Post("/categories", categoryController.Store)

		// Update existing category
		router.Put("/categories/{id}", categoryController.Update)

		// Delete category
		router.Delete("/categories/{id}", categoryController.Destroy)
	})

	// Article/Product management routes (methodes/admin only)
	articleController := controllers.NewArticleController()
	facades.Route().Middleware(middleware.Auth()).Group(func(router route.Router) {
		// List articles with pagination, search and filtering
		router.Get("/articles", articleController.Index)

		// Get specific article with all relationships
		router.Get("/articles/{id}", articleController.Show)

		// Create new article (Step 1: Basic Info)
		router.Post("/articles", articleController.Store)

		// Update existing article
		router.Put("/articles/{id}", articleController.Update)

		// Delete article
		router.Delete("/articles/{id}", articleController.Destroy)

		// Step 2: Attributes definition
		router.Post("/articles/{id}/attributes", articleController.CreateAttribute)
		router.Get("/articles/{id}/attributes", articleController.GetAttributes)

		// Step 3: Upload multiple images
		router.Post("/articles/{id}/images", articleController.CreateImages)
		router.Get("/articles/{id}/images", articleController.GetImages)

		// Step 4: Add product variants and set attribute values
		router.Post("/articles/{id}/variants", articleController.CreateVariant)
		router.Get("/articles/{id}/variants", articleController.GetVariants)

		// Note: Step 5 (Define storage location) is handled in the main article creation/update
		// Note: Step 6 (Define recipe) will be implemented separately as recipe management
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
