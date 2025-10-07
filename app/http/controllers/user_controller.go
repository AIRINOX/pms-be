package controllers

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"

	"pms/app/models"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

// CreateUserRequest represents the user creation request payload
type CreateUserRequest struct {
	Username string `json:"username" form:"username" validate:"required|min_len:3|max_len:50"`
	Password string `json:"password" form:"password" validate:"required|min_len:6"`
	Name     string `json:"name" form:"name" validate:"required|min_len:2|max_len:100"`
	Email    string `json:"email" form:"email" validate:"required|email"`
	Phone    string `json:"phone" form:"phone" validate:"max_len:20"`
	RoleID   uint   `json:"role_id" form:"role_id" validate:"required|numeric"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

// UpdateUserRequest represents the user update request payload
type UpdateUserRequest struct {
	Username string `json:"username" form:"username" validate:"min_len:3|max_len:50"`
	Password string `json:"password" form:"password" validate:"min_len:6"`
	Name     string `json:"name" form:"name" validate:"min_len:2|max_len:100"`
	Email    string `json:"email" form:"email" validate:"email"`
	Phone    string `json:"phone" form:"phone" validate:"max_len:20"`
	RoleID   uint   `json:"role_id" form:"role_id" validate:"numeric"`
	IsActive *bool  `json:"is_active" form:"is_active"`
}

// isAdmin checks if the authenticated user is an admin
func (r *UserController) isAdmin(ctx http.Context) bool {
	var user models.User
	err := facades.Auth(ctx).User(&user) // Get user directly into the model
	if err != nil {
		return false
	}

	// Load the role relationship
	if err := facades.Orm().Query().With("Role").Where("id", user.ID).FirstOrFail(&user); err != nil {
		return false
	}

	return user.Role.Key == "admin"
}

// Index returns a paginated list of users (admin only)
func (r *UserController) Index(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	// Parse new query parameters
	pageIndex, _ := strconv.Atoi(ctx.Request().Query("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "10"))
	searchQuery := ctx.Request().Query("query", "")
	sortKey := ctx.Request().Query("sort[key]", "")
	sortOrder := ctx.Request().Query("sort[order]", "")

	// Parse filter data
	filterName := ctx.Request().Query("filterData[name]", "")
	filterRoles := ctx.Request().QueryArray("filterData[role][]")
	filterStatuses := ctx.Request().QueryArray("filterData[status][]")

	fmt.Println(filterRoles)
	fmt.Println(filterStatuses)
	query := facades.Orm().Query().With("Role")

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("username LIKE ? OR name LIKE ? OR email LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Apply name filter
	if filterName != "" {
		query = query.Where("name LIKE ?", "%"+filterName+"%")
	}

	// Apply role filter
	if len(filterRoles) > 0 {
		var roleValues []any
		for _, role := range filterRoles {
			roleValues = append(roleValues, role)
		}
		query = query.Join("left join roles on users.role_id = roles.id")
		query = query.WhereIn("roles.key", roleValues)
	}

	// Apply status filter
	if len(filterStatuses) > 0 {
		var statusValues []any
		for _, status := range filterStatuses {
			switch status {
			case "true":
				statusValues = append(statusValues, true)
			case "false":
				statusValues = append(statusValues, false)
			}
		}
		if len(statusValues) > 0 {
			query = query.WhereIn("is_active", statusValues)
		}
	}

	// Apply sorting
	if sortKey != "" && sortOrder != "" {
		query = query.OrderBy(sortKey, sortOrder)
	} else {
		query = query.OrderBy("id", "desc") // default sorting
	}

	var users []models.User

	// Get total count
	total, err := query.Model(&models.User{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count users",
		})
	}

	// Get paginated results
	offset := (pageIndex - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve users",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"users": users,
		"pagination": http.Json{
			"current_page": pageIndex,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific user by ID (admin only)
func (r *UserController) Show(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "User ID is required",
		})
	}

	var user models.User
	if err := facades.Orm().Query().With("Role").Where("id", id).FirstOrFail(&user); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "User not found",
				"message": "The requested user does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve user",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"user": user,
	})
}

// Store creates a new user (admin only)
func (r *UserController) Store(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	var request CreateUserRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// log request content in console
	facades.Log().Info(request.Name)
	// not in file i need log in console cli

	// Validate the request data
	validator, err := facades.Validation().Make(map[string]any{
		"username":  request.Username,
		"password":  request.Password,
		"name":      request.Name,
		"email":     request.Email,
		"phone":     request.Phone,
		"role_id":   request.RoleID,
		"is_active": request.IsActive,
	}, map[string]string{
		"username": "required|min_len:3|max_len:50",
		"password": "required|min_len:6",
		"name":     "required|min_len:2|max_len:100",
		"email":    "required|email",
		"phone":    "max_len:20",
		"role_id":  "required|numeric",
	})

	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error": "Validation error",
		})
	}

	if validator.Fails() {
		return ctx.Response().Status(422).Json(http.Json{
			"error":  "Validation failed",
			"errors": validator.Errors().All(),
		})
	}

	// Check if username already exists
	var existingUser models.User
	if err := facades.Orm().Query().Where("username", request.Username).FirstOrFail(&existingUser); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Username already exists",
			"message": "A user with this username already exists",
		})
	}

	// Check if email already exists
	if err := facades.Orm().Query().Where("email", request.Email).FirstOrFail(&existingUser); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Email already exists",
			"message": "A user with this email already exists",
		})
	}

	// Verify role exists
	var role models.Role
	if err := facades.Orm().Query().Where("id", request.RoleID).FirstOrFail(&role); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid role",
			"message": "The specified role does not exist",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to hash password",
			"message": "Internal server error",
		})
	}

	// Set default value for IsActive if not provided
	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	// Create new user
	user := models.User{
		Username: request.Username,
		Password: string(hashedPassword),
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		RoleID:   request.RoleID,
		IsActive: isActive,
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to create user",
			"message": "Internal server error",
		})
	}

	// Load the role relationship
	facades.Orm().Query().With("Role").Where("id", user.ID).First(&user)

	return ctx.Response().Status(201).Json(http.Json{
		"message": "User created successfully",
		"user":    user,
	})
}

// Update modifies an existing user (admin only)
func (r *UserController) Update(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "User ID is required",
		})
	}

	var request UpdateUserRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing user
	var user models.User
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&user); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "User not found",
				"message": "The requested user does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve user",
		})
	}

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Username != "" {
		validationData["username"] = request.Username
		validationRules["username"] = "min_len:3|max_len:50"

		// Check if username already exists (excluding current user)
		var existingUser models.User
		if err := facades.Orm().Query().Where("username", request.Username).Where("id != ?", id).FirstOrFail(&existingUser); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Username already exists",
				"message": "A user with this username already exists",
			})
		}
	}

	if request.Email != "" {
		validationData["email"] = request.Email
		validationRules["email"] = "email"

		// Check if email already exists (excluding current user)
		var existingUser models.User
		if err := facades.Orm().Query().Where("email", request.Email).Where("id != ?", id).FirstOrFail(&existingUser); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Email already exists",
				"message": "A user with this email already exists",
			})
		}
	}

	if request.Password != "" {
		validationData["password"] = request.Password
		validationRules["password"] = "min_len:6"
	}

	if request.Name != "" {
		validationData["name"] = request.Name
		validationRules["name"] = "min_len:2|max_len:100"
	}

	if request.Phone != "" {
		validationData["phone"] = request.Phone
		validationRules["phone"] = "max_len:20"
	}

	if request.RoleID != 0 {
		validationData["role_id"] = request.RoleID
		validationRules["role_id"] = "numeric"

		// Verify role exists
		var role models.Role
		if err := facades.Orm().Query().Where("id", request.RoleID).FirstOrFail(&role); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid role",
				"message": "The specified role does not exist",
			})
		}
	}

	// Validate the request data
	if len(validationData) > 0 {
		validator, err := facades.Validation().Make(validationData, validationRules)
		if err != nil {
			return ctx.Response().Status(500).Json(http.Json{
				"error": "Validation error",
			})
		}

		if validator.Fails() {
			return ctx.Response().Status(422).Json(http.Json{
				"error":  "Validation failed",
				"errors": validator.Errors().All(),
			})
		}
	}

	// Update user fields
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.RoleID != 0 {
		user.RoleID = request.RoleID
	}
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}

	// Hash new password if provided
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Failed to hash password",
				"message": "Internal server error",
			})
		}
		user.Password = string(hashedPassword)
	}

	// Save updated user
	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to update user",
			"message": "Internal server error",
		})
	}

	// Load the role relationship
	facades.Orm().Query().With("Role").Where("id", user.ID).First(&user)

	return ctx.Response().Status(200).Json(http.Json{
		"message": "User updated successfully",
		"user":    user,
	})
}

// Destroy deletes a user (admin only)
func (r *UserController) Destroy(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "User ID is required",
		})
	}

	// Find existing user
	var user models.User
	if err := facades.Orm().Query().With("Role").Where("id", id).FirstOrFail(&user); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "User not found",
				"message": "The requested user does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve user",
		})
	}

	// Prevent admin from deleting themselves

	var currentUser models.User
	err := facades.Auth(ctx).User(&currentUser) // Get user directly into the model
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to retrieve user",
			"message": "Internal server error",
		})
	}

	if currentUser.ID == user.ID {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete yourself",
			"message": "You cannot delete your own account",
		})
	}

	// Delete the user
	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to delete user",
			"message": "Internal server error",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "User deleted successfully",
	})
}

// ToggleStatus toggles user active status (admin only)
func (r *UserController) ToggleStatus(ctx http.Context) http.Response {
	if !r.isAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "User ID is required",
		})
	}

	// Find existing user
	var user models.User
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&user); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "User not found",
				"message": "The requested user does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve user",
		})
	}

	// Prevent admin from deactivating themselves
	var currentUser models.User
	err := facades.Auth(ctx).User(&currentUser) // Get user directly into the model
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to retrieve user",
			"message": "Internal server error",
		})
	}
	if currentUser.ID == user.ID && user.IsActive {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot deactivate yourself",
			"message": "You cannot deactivate your own account",
		})
	}

	// Toggle status
	user.IsActive = !user.IsActive

	// Save updated user
	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to update user status",
			"message": "Internal server error",
		})
	}

	// Load the role relationship
	facades.Orm().Query().With("Role").Where("id", user.ID).First(&user)

	status := "activated"
	if !user.IsActive {
		status = "deactivated"
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "User " + status + " successfully",
		"user":    user,
	})
}
