package middleware

import (
	"errors"
	"pms/app/models"

	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		// Get the token from the request
		token := ctx.Request().Header("Authorization")
		if token == "" {
			ctx.Response().Status(401).Json(http.Json{
				"error":   "Unauthorized",
				"message": "Authorization token required",
			})
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Parse the token and check for expiration
		payload, err := facades.Auth(ctx).Parse(token)
		if err != nil {
			if errors.Is(err, auth.ErrorTokenExpired) {
				ctx.Response().Status(401).Json(http.Json{
					"error":   "Token Expired",
					"message": "Your authentication token has expired",
				})
				return
			}
			ctx.Response().Status(401).Json(http.Json{
				"error":   "Unauthorized",
				"message": "Invalid authentication token",
			})
			return
		}

		// Optionally, you can still get the user if needed
		var user models.User
		err = facades.Auth(ctx).User(&user)
		if err != nil {
			ctx.Response().Status(401).Json(http.Json{
				"error":   "Unauthorized",
				"message": "User not found",
			})
			return
		}

		// Store payload in context for later use if needed
		ctx.WithValue("auth_payload", payload)
		ctx.WithValue("user", user)
	}
}

func Guest() http.Middleware {
	return func(ctx http.Context) {
		// Get the token from the request
		token := ctx.Request().Header("Authorization")
		if token == "" {
			// No token, user is guest - continue
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Parse the token to check if user is authenticated
		_, err := facades.Auth(ctx).Parse(token)
		if err == nil {
			// Token is valid, user is authenticated
			ctx.Response().Status(400).Json(http.Json{
				"error":   "Already authenticated",
				"message": "You are already logged in",
			})
			return
		}

		// Token is invalid or expired, user is guest - continue
	}
}
