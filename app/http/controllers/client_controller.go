package controllers

import (
	"strconv"

	"github.com/goravel/framework/errors"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type ClientController struct {
	// Dependent services
}

func NewClientController() *ClientController {
	return &ClientController{
		// Inject services
	}
}

// CreateClientRequest represents the client creation request payload
type CreateClientRequest struct {
	Name    string `json:"name" form:"name" validate:"required|min_len:2|max_len:255"`
	Phone   string `json:"phone" form:"phone" validate:"max_len:20"`
	Email   string `json:"email" form:"email" validate:"email|max_len:255"`
	Address string `json:"address" form:"address"`
}

// UpdateClientRequest represents the client update request payload
type UpdateClientRequest struct {
	Name    string `json:"name" form:"name" validate:"min_len:2|max_len:255"`
	Phone   string `json:"phone" form:"phone" validate:"max_len:20"`
	Email   string `json:"email" form:"email" validate:"email|max_len:255"`
	Address string `json:"address" form:"address"`
}

// isCommercialOrAdmin checks if the authenticated user is commercial or admin
func (r *ClientController) isCommercialOrAdmin(ctx http.Context) bool {
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

// Index returns a paginated list of clients with search and filtering
func (r *ClientController) Index(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	// Parse query parameters
	pageIndex, _ := strconv.Atoi(ctx.Request().Query("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "10"))
	searchQuery := ctx.Request().Query("query", "")
	sortKey := ctx.Request().Query("sort[key]", "name")
	sortOrder := ctx.Request().Query("sort[order]", "asc")

	// Parse filter data
	filterName := ctx.Request().Query("filterData[name]", "")
	filterEmail := ctx.Request().Query("filterData[email]", "")
	filterPhone := ctx.Request().Query("filterData[phone]", "")

	query := facades.Orm().Query().With("ClientSites")

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ? OR address LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Apply specific filters
	if filterName != "" {
		query = query.Where("name LIKE ?", "%"+filterName+"%")
	}
	if filterEmail != "" {
		query = query.Where("email LIKE ?", "%"+filterEmail+"%")
	}
	if filterPhone != "" {
		query = query.Where("phone LIKE ?", "%"+filterPhone+"%")
	}

	// Apply sorting
	if sortKey != "" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.OrderBy(sortKey, sortOrder)
	} else {
		query = query.OrderBy("name", "asc")
	}

	// Get total count for pagination
	total, err := query.Model(&models.Client{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count clients",
		})
	}

	// Apply pagination
	offset := (pageIndex - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	var clients []models.Client
	if err := query.Find(&clients); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve clients",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"clients": clients,
		"pagination": http.Json{
			"current_page": pageIndex,
			"per_page":     pageSize,
			"total":        total,
			"last_page":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific client by ID with its sites and orders
func (r *ClientController) Show(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	id := ctx.Request().Route("clientId")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	var client models.Client
	if err := facades.Orm().Query().With("ClientSites").Where("id", id).First(&client); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
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

	return ctx.Response().Status(200).Json(http.Json{
		"client": client,
	})
}

// Store creates a new client
func (r *ClientController) Store(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	var request CreateClientRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate the request data
	validator, err := facades.Validation().Make(map[string]any{
		"name":    request.Name,
		"phone":   request.Phone,
		"email":   request.Email,
		"address": request.Address,
	}, map[string]string{
		"name":  "required|min_len:2|max_len:255",
		"phone": "max_len:20",
		"email": "email|max_len:255",
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

	// Check if client name already exists
	if request.Name != "" {
		var existingClient models.Client
		if err := facades.Orm().Query().Where("name", request.Name).FirstOrFail(&existingClient); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Client name already exists",
				"message": "A client with this name already exists",
			})
		}
	}

	// Check if email already exists (if provided)
	if request.Email != "" {
		var existingClient models.Client
		if err := facades.Orm().Query().Where("email", request.Email).FirstOrFail(&existingClient); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Email already exists",
				"message": "A client with this email already exists",
			})
		}
	}

	// Create new client
	client := models.Client{
		Name:    request.Name,
		Phone:   request.Phone,
		Email:   request.Email,
		Address: request.Address,
	}

	if err := facades.Orm().Query().Create(&client); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to create client",
			"message": "Internal server error",
		})
	}

	// Load relationships
	facades.Orm().Query().With("ClientSites").Where("id", client.ID).First(&client)

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Client created successfully",
		"client":  client,
	})
}

// Update modifies an existing client
func (r *ClientController) Update(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	id := ctx.Request().Route("clientId")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	var request UpdateClientRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing client
	var client models.Client
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&client); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
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

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Name != "" {
		validationData["name"] = request.Name
		validationRules["name"] = "min_len:2|max_len:255"

		// Check if name already exists (excluding current client)
		var existingClient models.Client
		if err := facades.Orm().Query().Where("name", request.Name).Where("id != ?", id).FirstOrFail(&existingClient); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Client name already exists",
				"message": "A client with this name already exists",
			})
		}
	}

	if request.Email != "" {
		validationData["email"] = request.Email
		validationRules["email"] = "email|max_len:255"

		// Check if email already exists (excluding current client)
		var existingClient models.Client
		if err := facades.Orm().Query().Where("email", request.Email).Where("id != ?", id).FirstOrFail(&existingClient); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Email already exists",
				"message": "A client with this email already exists",
			})
		}
	}

	if request.Phone != "" {
		validationData["phone"] = request.Phone
		validationRules["phone"] = "max_len:20"
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

	// Update client fields
	if request.Name != "" {
		client.Name = request.Name
	}
	if request.Phone != "" {
		client.Phone = request.Phone
	}
	if request.Email != "" {
		client.Email = request.Email
	}
	if request.Address != "" {
		client.Address = request.Address
	}

	// Save updated client
	if err := facades.Orm().Query().Save(&client); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to update client",
			"message": "Internal server error",
		})
	}

	// Load relationships
	facades.Orm().Query().With("ClientSites").Where("id", client.ID).First(&client)

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Client updated successfully",
		"client":  client,
	})
}

// Destroy deletes a client
func (r *ClientController) Destroy(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	id := ctx.Request().Route("clientId")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	// Find existing client
	var client models.Client
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&client); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
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

	// Check if client has any orders
	orderCount, err := facades.Orm().Query().Model(&models.OrderFabrication{}).Where("client_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check client orders",
		})
	}
	if orderCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete client",
			"message": "Client has existing fabrication orders and cannot be deleted",
		})
	}

	// Delete the client (this will cascade delete client sites due to foreign key)
	if _, err := facades.Orm().Query().Delete(&client); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Failed to delete client",
			"message": "Internal server error",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Client deleted successfully",
	})
}

// GetClientOrders returns all orders for a specific client
func (r *ClientController) GetClientOrders(ctx http.Context) http.Response {
	if !r.isCommercialOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Commercial or Admin access required",
		})
	}

	id := ctx.Request().Route("clientId")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Client ID is required",
		})
	}

	// Verify client exists
	var client models.Client
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&client); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
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

	// Parse pagination parameters
	pageIndex, _ := strconv.Atoi(ctx.Request().Query("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(ctx.Request().Query("pageSize", "10"))
	status := ctx.Request().Query("status", "")

	query := facades.Orm().Query().With("Article", "Variant", "ClientSite", "Creator").Where("client_id", id)

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status", status)
	}

	// Get total count
	total, err := query.Model(&models.OrderFabrication{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count client orders",
		})
	}

	// Apply pagination and ordering
	offset := (pageIndex - 1) * pageSize
	query = query.OrderBy("created_at", "desc").Offset(offset).Limit(pageSize)

	var orders []models.OrderFabrication
	if err := query.Find(&orders); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve client orders",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"client": client,
		"orders": orders,
		"pagination": http.Json{
			"current_page": pageIndex,
			"per_page":     pageSize,
			"total":        total,
			"last_page":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}
