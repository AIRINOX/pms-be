package middleware

import (
	"pms/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		// Check if user is authenticated using JWT
		var user models.User
		err := facades.Auth(ctx).User(&user) // Get user directly into the model
		if err != nil {
			ctx.Response().Status(401).Json(http.Json{
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			return
		}
	}
}

func Guest() http.Middleware {
	return func(ctx http.Context) {
		// Check if user is already authenticated
		var user models.User
		err := facades.Auth(ctx).User(&user) // Get user directly into the model
		if err == nil {
			ctx.Response().Status(400).Json(http.Json{
				"error":   "Already authenticated",
				"message": "You are already logged in",
			})
			return
		}
	}
}
