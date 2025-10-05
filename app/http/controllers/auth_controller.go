package controllers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pms/app/models"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// Login handles user authentication
func (r *AuthController) Login(ctx http.Context) http.Response {
	var request LoginRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate the request data
	validator, err := facades.Validation().Make(map[string]any{
		"username": request.Username,
		"password": request.Password,
	}, map[string]string{
		"username": "required",
		"password": "required",
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

	// Find user by username
	var user models.User
	if err := facades.Orm().Query().Where("username", request.Username).First(&user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(401).Json(http.Json{
				"error":   "Invalid credentials",
				"message": "Username or password is incorrect",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Internal server error",
		})
	}

	// Check if user is active
	if !user.IsActive {
		return ctx.Response().Status(401).Json(http.Json{
			"error":   "Account disabled",
			"message": "Your account has been disabled",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return ctx.Response().Status(401).Json(http.Json{
			"error":   "Invalid credentials",
			"message": "Username or password is incorrect",
		})
	}

	// Generate JWT token
	token, err := facades.Auth().LoginUsingID(user.ID)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to generate token",
			"message": "Internal server error",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Login successful",
		"user": http.Json{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"phone":    user.Phone,
			"role_id":  user.RoleID,
			"is_active": user.IsActive,
		},
		"token": token,
	})
}

// Logout handles user logout
func (r *AuthController) Logout(ctx http.Context) http.Response {
	if err := facades.Auth().Logout(); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to logout",
			"message": "Internal server error",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Logout successful",
	})
}

// Me returns the authenticated user's information
func (r *AuthController) Me(ctx http.Context) http.Response {
	user := facades.Auth().User(ctx)
	if user == nil {
		return ctx.Response().Status(401).Json(http.Json{
			"error":   "Unauthorized",
			"message": "User not authenticated",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"user": user,
	})
}

// RefreshToken refreshes the JWT token
func (r *AuthController) RefreshToken(ctx http.Context) http.Response {
	token, err := facades.Auth().Refresh()
	if err != nil {
		return ctx.Response().Status(401).Json(http.Json{
			"error":   "Failed to refresh token",
			"message": "Invalid or expired token",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Token refreshed successfully",
		"token":   token,
	})
}