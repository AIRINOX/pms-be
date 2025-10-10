package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type StorageLocationController struct {
	// Dependent services
}

func NewStorageLocationController() *StorageLocationController {
	return &StorageLocationController{
		// Inject services
	}
}

// CreateStorageLocationRequest represents the storage location creation request payload
type CreateStorageLocationRequest struct {
	Name        string `json:"name" form:"name" validate:"required|min_len:2|max_len:255"`
	Description string `json:"description" form:"description"`
}

// UpdateStorageLocationRequest represents the storage location update request payload
type UpdateStorageLocationRequest struct {
	Name        string `json:"name" form:"name" validate:"min_len:2|max_len:255"`
	Description string `json:"description" form:"description"`
}

// isMethodesOrAdmin checks if the authenticated user is methodes or admin
func (r *StorageLocationController) isMethodesOrAdmin(ctx http.Context) bool {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return false
	}

	// Load the role relationship
	if err := facades.Orm().Query().With("Role").Where("id", user.ID).First(&user); err != nil {
		return false
	}

	return user.Role.Key == "admin" || user.Role.Key == "methodes"
}

// Index returns a paginated list of storage locations with search and filtering
func (r *StorageLocationController) Index(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
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

	query := facades.Orm().Query()

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Apply specific filters
	if filterName != "" {
		query = query.Where("name LIKE ?", "%"+filterName+"%")
	}

	// Apply sorting
	if sortKey != "" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.OrderBy(sortKey, sortOrder)
	} else {
		query = query.OrderBy("name", "asc")
	}

	var storageLocations []models.StorageLocation

	// Get total count
	total, err := query.Model(&models.StorageLocation{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count storage locations",
		})
	}

	// Get paginated results
	offset := (pageIndex - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&storageLocations); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve storage locations",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"storage_locations": storageLocations,
		"pagination": http.Json{
			"current_page": pageIndex,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific storage location by ID
func (r *StorageLocationController) Show(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Storage location ID is required",
		})
	}

	var storageLocation models.StorageLocation
	if err := facades.Orm().Query().Where("id", id).First(&storageLocation); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Storage location not found",
				"message": "The requested storage location does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve storage location",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"storage_location": storageLocation,
	})
}

// Store creates a new storage location
func (r *StorageLocationController) Store(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	var request CreateStorageLocationRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate input
	validator, err := facades.Validation().Make(map[string]any{
		"name":        request.Name,
		"description": request.Description,
	}, map[string]string{
		"name":        "required|min_len:2|max_len:255",
		"description": "max_len:1000",
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

	// Check if name already exists
	var existingLocation models.StorageLocation
	if err := facades.Orm().Query().Where("name", request.Name).FirstOrFail(&existingLocation); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Storage location name already exists",
			"message": "A storage location with this name already exists",
		})
	}

	// Create new storage location
	storageLocation := models.StorageLocation{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := facades.Orm().Query().Create(&storageLocation); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create storage location",
		})
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message":          "Storage location created successfully",
		"storage_location": storageLocation,
	})
}

// Update modifies an existing storage location
func (r *StorageLocationController) Update(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Storage location ID is required",
		})
	}

	var request UpdateStorageLocationRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing storage location
	var storageLocation models.StorageLocation
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&storageLocation); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Storage location not found",
				"message": "The requested storage location does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve storage location",
		})
	}

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Name != "" {
		validationData["name"] = request.Name
		validationRules["name"] = "min_len:2|max_len:255"

		// Check if name already exists (excluding current storage location)
		var existingLocation models.StorageLocation
		if err := facades.Orm().Query().Where("name", request.Name).Where("id != ?", id).FirstOrFail(&existingLocation); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Storage location name already exists",
				"message": "A storage location with this name already exists",
			})
		}
	}

	if request.Description != "" {
		validationData["description"] = request.Description
		validationRules["description"] = "max_len:1000"
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

	// Update storage location fields
	if request.Name != "" {
		storageLocation.Name = request.Name
	}
	if request.Description != "" {
		storageLocation.Description = request.Description
	}

	// Save changes
	if err := facades.Orm().Query().Save(&storageLocation); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to update storage location",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message":          "Storage location updated successfully",
		"storage_location": storageLocation,
	})
}

// Destroy deletes a storage location
func (r *StorageLocationController) Destroy(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	id := ctx.Request().Route("id")
	if id == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Storage location ID is required",
		})
	}

	// Find existing storage location
	var storageLocation models.StorageLocation
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&storageLocation); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Storage location not found",
				"message": "The requested storage location does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve storage location",
		})
	}

	// Check if storage location has any products
	productCount, err := facades.Orm().Query().Model(&models.Product{}).Where("location_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check storage location usage",
		})
	}
	if productCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete storage location",
			"message": "Storage location has existing products and cannot be deleted",
		})
	}

	// Check if storage location has any stock levels
	stockLevelCount, err := facades.Orm().Query().Model(&models.StockLevel{}).Where("location_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check storage location stock levels",
		})
	}
	if stockLevelCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete storage location",
			"message": "Storage location has existing stock levels and cannot be deleted",
		})
	}

	// Check if storage location has any stock movements
	stockMovementCount, err := facades.Orm().Query().Model(&models.StockMovement{}).Where("location_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check storage location stock movements",
		})
	}
	if stockMovementCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete storage location",
			"message": "Storage location has existing stock movements and cannot be deleted",
		})
	}

	// Delete storage location
	if _, err := facades.Orm().Query().Delete(&storageLocation); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete storage location",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Storage location deleted successfully",
	})
}

func (r *StorageLocationController) Select(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	searchQuery := ctx.Request().Query("query", "")

	query := facades.Orm().Query()
	if searchQuery != "" {
		query = query.Where("name LIKE ?", "%"+searchQuery+"%")
	}

	var storageLocations []models.StorageLocation
	if err := query.OrderBy("name", "asc").Limit(20).Find(&storageLocations); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve storage location options",
		})
	}

	options := make([]http.Json, 0, len(storageLocations))
	for _, s := range storageLocations {
		options = append(options, http.Json{
			"value": s.ID,
			"label": s.Name,
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"options": options,
	})
}
