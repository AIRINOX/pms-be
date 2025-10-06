package controllers

import (
	"errors"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"gorm.io/gorm"

	"pms/app/models"
)

type ClientSiteController struct {
	// Dependent services
}

func NewClientSiteController() *ClientSiteController {
	return &ClientSiteController{
		// Inject services
	}
}

// CreateClientSiteRequest represents the client site creation request payload
type CreateClientSiteRequest struct {
	ClientID     uint   `json:"client_id" form:"client_id" validate:"required|numeric"`
	Title        string `json:"title" form:"title" validate:"required|min_len:2|max_len:255"`
	Address      string `json:"address" form:"address"`
	ContactName  string `json:"contact_name" form:"contact_name" validate:"max_len:255"`
	ContactPhone string `json:"contact_phone" form:"contact_phone" validate:"max_len:20"`
}

// UpdateClientSiteRequest represents the client site update request payload
type UpdateClientSiteRequest struct {
	Title        string `json:"title" form:"title" validate:"min_len:2|max_len:255"`
	Address      string `json:"address" form:"address"`
	ContactName  string `json:"contact_name" form:"contact_name" validate:"max_len:255"`
	ContactPhone string `json:"contact_phone" form:"contact_phone" validate:"max_len:20"`
}

// isCommercialOrAdmin checks if the authenticated user is commercial or admin
func (r *ClientSiteController) isCommercialOrAdmin(ctx http.Context) bool {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return false
	}

	// Load the role relationship
	if err := facades.Orm().Query().With("Role").Where("id", user.ID).First(&user); err != nil {
		return false
	}

	return user.Role.Key == "admin" || user.Role.Key == "commercial"
}

// Index returns all sites for a specific client
func (r *ClientSiteController) Index(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	clientID := ctx.Request().Route("clientId")
	if clientID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	// Verify client exists
	var client models.Client
	if err := facades.Orm().Query().Where("id", clientID).FirstOrFail(&client); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Client not found",
				"message": "The requested client does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client",
		})
	}

	// Parse query parameters
	pageIndex, _ := strconv.Atoi(ctx.Request().Query("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "10"))
	searchQuery := ctx.Request().Query("query", "")

	query := facades.Orm().Query().With("Client").Where("client_id", clientID)

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("title LIKE ? OR address LIKE ? OR contact_name LIKE ? OR contact_phone LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Get total count
	total, err := query.Model(&models.ClientSite{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count client sites",
		})
	}

	// Apply pagination and ordering
	offset := (pageIndex - 1) * pageSize
	query = query.OrderBy("title", "asc").Offset(offset).Limit(pageSize)

	var sites []models.ClientSite
	if err := query.Find(&sites); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client sites",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"client": client,
		"sites":  sites,
		"pagination": http.Json{
			"current_page": pageIndex,
			"per_page":     pageSize,
			"total":        total,
			"last_page":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific client site by ID
func (r *ClientSiteController) Show(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	clientID := ctx.Request().Route("clientId")
	siteID := ctx.Request().Route("siteId")
	
	if clientID == "" || siteID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID and Site ID are required",
		})
	}

	var site models.ClientSite
	if err := facades.Orm().Query().With("Client").Where("id", siteID).Where("client_id", clientID).First(&site); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Client site not found",
				"message": "The requested client site does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client site",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"site": site,
	})
}

// Store creates a new client site
func (r *ClientSiteController) Store(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	clientID := ctx.Request().Route("clientId")
	if clientID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	var request CreateClientSiteRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Set client ID from route
	clientIDUint, err := strconv.ParseUint(clientID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid client ID",
			"message": "Client ID must be a valid number",
		})
	}
	request.ClientID = uint(clientIDUint)

	// Verify client exists
	var client models.Client
	if err := facades.Orm().Query().Where("id", request.ClientID).FirstOrFail(&client); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Client not found",
				"message": "The specified client does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client",
		})
	}

	// Validate the request data
	validator, err := facades.Validation().Make(map[string]any{
		"client_id":     request.ClientID,
		"title":         request.Title,
		"address":       request.Address,
		"contact_name":  request.ContactName,
		"contact_phone": request.ContactPhone,
	}, map[string]string{
		"client_id":     "required|numeric",
		"title":         "required|min_len:2|max_len:255",
		"contact_name":  "max_len:255",
		"contact_phone": "max_len:20",
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

	// Check if site title already exists for this client
	var existingSite models.ClientSite
	if err := facades.Orm().Query().Where("client_id", request.ClientID).Where("title", request.Title).FirstOrFail(&existingSite); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Site title already exists",
			"message": "A site with this title already exists for this client",
		})
	}

	// Create new client site
	site := models.ClientSite{
		ClientID:     request.ClientID,
		Title:        request.Title,
		Address:      request.Address,
		ContactName:  request.ContactName,
		ContactPhone: request.ContactPhone,
	}

	if err := facades.Orm().Query().Create(&site); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to create client site",
			"message": "Internal server error",
		})
	}

	// Load relationships
	facades.Orm().Query().With("Client").Where("id", site.ID).First(&site)

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Client site created successfully",
		"site":    site,
	})
}

// Update modifies an existing client site
func (r *ClientSiteController) Update(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	clientID := ctx.Request().Route("clientId")
	siteID := ctx.Request().Route("siteId")
	
	if clientID == "" || siteID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID and Site ID are required",
		})
	}

	var request UpdateClientSiteRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing client site
	var site models.ClientSite
	if err := facades.Orm().Query().Where("id", siteID).Where("client_id", clientID).FirstOrFail(&site); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Client site not found",
				"message": "The requested client site does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client site",
		})
	}

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Title != "" {
		validationData["title"] = request.Title
		validationRules["title"] = "min_len:2|max_len:255"

		// Check if title already exists for this client (excluding current site)
		var existingSite models.ClientSite
		if err := facades.Orm().Query().Where("client_id", clientID).Where("title", request.Title).Where("id != ?", siteID).FirstOrFail(&existingSite); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Site title already exists",
				"message": "A site with this title already exists for this client",
			})
		}
	}

	if request.ContactName != "" {
		validationData["contact_name"] = request.ContactName
		validationRules["contact_name"] = "max_len:255"
	}

	if request.ContactPhone != "" {
		validationData["contact_phone"] = request.ContactPhone
		validationRules["contact_phone"] = "max_len:20"
	}

	// Validate if there's data to validate
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

	// Update site fields
	if request.Title != "" {
		site.Title = request.Title
	}
	if request.Address != "" {
		site.Address = request.Address
	}
	if request.ContactName != "" {
		site.ContactName = request.ContactName
	}
	if request.ContactPhone != "" {
		site.ContactPhone = request.ContactPhone
	}

	// Save updated site
	if err := facades.Orm().Query().Save(&site); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to update client site",
			"message": "Internal server error",
		})
	}

	// Load relationships
	facades.Orm().Query().With("Client").Where("id", site.ID).First(&site)

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Client site updated successfully",
		"site":    site,
	})
}

// Destroy deletes a client site
func (r *ClientSiteController) Destroy(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	clientID := ctx.Request().Route("clientId")
	siteID := ctx.Request().Route("siteId")
	
	if clientID == "" || siteID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID and Site ID are required",
		})
	}

	// Find existing client site
	var site models.ClientSite
	if err := facades.Orm().Query().Where("id", siteID).Where("client_id", clientID).FirstOrFail(&site); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Client site not found",
				"message": "The requested client site does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client site",
		})
	}

	// Check if site has any orders
	orderCount, err := facades.Orm().Query().Model(&models.OrderFabrication{}).Where("client_site_id", siteID).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check site orders",
		})
	}
	if orderCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete client site",
			"message": "Client site has existing fabrication orders and cannot be deleted",
		})
	}

	// Delete the client site
	if _, err := facades.Orm().Query().Delete(&site); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to delete client site",
			"message": "Internal server error",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Client site deleted successfully",
	})
}