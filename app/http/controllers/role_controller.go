package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type RoleController struct {
	// Dependent services
}

func NewRoleController() *RoleController {
	return &RoleController{
		// Inject services
	}
}

// isAdmin checks if the authenticated user is an admin
func (r *RoleController) isAdmin(ctx http.Context) bool {
	var user models.User
	err := facades.Auth(ctx).User(&user) // Get user directly into the model
	if err != nil {
		return false
	}

	// Load the role relationship
	if err := facades.Orm().Query().With("Role").Where("id", user.ID).First(&user); err != nil {
		return false
	}

	return user.Role.Key == "admin"
}

// Index returns all roles (admin only)
func (r *RoleController) Index(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	var roles []models.Role
	if err := facades.Orm().Query().OrderBy("order_index").Find(&roles); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve roles",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"roles": roles,
	})
}
