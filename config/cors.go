package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("cors", map[string]any{
		// Cross-Origin Resource Sharing (CORS) Configuration
		//
		// Here you may configure your settings for cross-origin resource sharing
		// or "CORS". This determines what cross-origin operations may execute
		// in web browsers. You are free to adjust these settings as needed.
		//
		// To learn more: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		"paths":                []string{"*"}, // Allow CORS for all paths
		"allowed_methods":      []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		"allowed_origins":      []string{"*"}, // Allow all origins - customize this in production
		"allowed_headers":      []string{"Content-Type", "X-CSRF-Token", "X-Requested-With", "Accept", "Origin", "Authorization"},
		"exposed_headers":      []string{"Content-Length", "Content-Type"},
		"max_age":              86400, // 24 hours
		"supports_credentials": true,  // Enable if you need to send cookies or auth headers
	})
}
