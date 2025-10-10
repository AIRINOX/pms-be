package controllers

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type CategoryController struct {
	// Dependent services
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		// Inject services
	}
}

// CreateCategoryRequest represents the category creation request payload
type CreateCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"required|min_len:2|max_len:255"`
	Description string `json:"description" form:"description"`
	ParentID    *uint  `json:"parent_id" form:"parent_id"`
}

// UpdateCategoryRequest represents the category update request payload
type UpdateCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"min_len:2|max_len:255"`
	Description string `json:"description" form:"description"`
	ParentID    *uint  `json:"parent_id" form:"parent_id"`
}

// isMethodesOrAdmin checks if the authenticated user is methodes or admin
func (r *CategoryController) isMethodesOrAdmin(ctx http.Context) bool {
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

// Index returns a paginated list of categories with search and filtering
func (r *CategoryController) Index(ctx http.Context) http.Response {
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
	sortKey := ctx.Request().Query("sort[key]", "title")
	sortOrder := ctx.Request().Query("sort[order]", "asc")
	includeChildren := ctx.Request().Query("include_children", "false")

	// Parse filter data
	filterTitle := ctx.Request().Query("filterData[title]", "")
	parentID := ctx.Request().Query("parent_id", "")

	query := facades.Orm().Query()

	// Include relationships if requested
	if includeChildren == "true" {
		query = query.With("Parent", "Children")
	} else {
		query = query.With("Parent")
	}

	// Apply parent filter
	if parentID != "" {
		if parentID == "null" {
			query = query.WhereNull("parent_id")
		} else {
			query = query.Where("parent_id", parentID)
		}
	}

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Apply specific filters
	if filterTitle != "" {
		query = query.Where("title LIKE ?", "%"+filterTitle+"%")
	}

	// Apply sorting
	if sortKey != "" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.OrderBy(sortKey, sortOrder)
	} else {
		query = query.OrderBy("title", "asc")
	}

	var categories []models.Category

	// Get total count
	total, err := query.Model(&models.Category{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count categories",
		})
	}

	// Get paginated results
	offset := (pageIndex - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&categories); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve categories",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"categories": categories,
		"pagination": http.Json{
			"current_page": pageIndex,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetTree returns categories in a hierarchical tree structure
func (r *CategoryController) GetTree(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	var categories []models.Category
	if err := facades.Orm().Query().With("Children").WhereNull("parent_id").OrderBy("title", "asc").Find(&categories); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve categories tree",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"categories": categories,
	})
}

// Show returns a specific category by ID
func (r *CategoryController) Show(ctx http.Context) http.Response {
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
			"message": "Category ID is required",
		})
	}

	var category models.Category
	if err := facades.Orm().Query().With("Parent").With("Children").Where("id", id).FirstOrFail(&category); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Category not found",
				"message": "The requested category does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve category",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"category": category,
	})
}

// Store creates a new category
func (r *CategoryController) Store(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	var request CreateCategoryRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate input
	validationData := map[string]any{
		"title":       request.Title,
		"description": request.Description,
	}
	validationRules := map[string]string{
		"title":       "required|min_len:2|max_len:255",
		"description": "max_len:1000",
	}

	if request.ParentID != nil {
		validationData["parent_id"] = *request.ParentID
		validationRules["parent_id"] = "numeric"
	}

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

	// Check if title already exists within the same parent
	query := facades.Orm().Query().Where("title", request.Title)
	if request.ParentID != nil {
		query = query.Where("parent_id", *request.ParentID)
	} else {
		query = query.WhereNull("parent_id")
	}

	var existingCategory models.Category
	if err := query.FirstOrFail(&existingCategory); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Category title already exists",
			"message": "A category with this title already exists in the same parent category",
		})
	}

	// Verify parent category exists if provided
	if request.ParentID != nil {
		var parentCategory models.Category
		if err := facades.Orm().Query().Where("id", *request.ParentID).FirstOrFail(&parentCategory); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid parent category",
				"message": "The specified parent category does not exist",
			})
		}
	}

	// Create new category
	category := models.Category{
		Title:       request.Title,
		Description: request.Description,
		ParentID:    request.ParentID,
	}

	if err := facades.Orm().Query().Create(&category); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create category",
		})
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Parent").Where("id", category.ID).First(&category); err != nil {
		// Category was created but failed to load relationships, still return success
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message":  "Category created successfully",
		"category": category,
	})
}

// Update modifies an existing category
func (r *CategoryController) Update(ctx http.Context) http.Response {
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
			"message": "Category ID is required",
		})
	}

	var request UpdateCategoryRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing category
	var category models.Category
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&category); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Category not found",
				"message": "The requested category does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve category",
		})
	}

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Title != "" {
		validationData["title"] = request.Title
		validationRules["title"] = "min_len:2|max_len:255"

		// Check if title already exists within the same parent (excluding current category)
		query := facades.Orm().Query().Where("title", request.Title).Where("id != ?", id)
		if request.ParentID != nil {
			query = query.Where("parent_id", *request.ParentID)
		} else {
			query = query.WhereNull("parent_id")
		}

		var existingCategory models.Category
		if err := query.FirstOrFail(&existingCategory); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "Category title already exists",
				"message": "A category with this title already exists in the same parent category",
			})
		}
	}

	if request.Description != "" {
		validationData["description"] = request.Description
		validationRules["description"] = "max_len:1000"
	}

	if request.ParentID != nil {
		validationData["parent_id"] = *request.ParentID
		validationRules["parent_id"] = "numeric"

		// Prevent setting parent to itself or its descendants
		if *request.ParentID == category.ID {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid parent category",
				"message": "A category cannot be its own parent",
			})
		}

		// Verify parent category exists
		var parentCategory models.Category
		if err := facades.Orm().Query().Where("id", *request.ParentID).FirstOrFail(&parentCategory); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid parent category",
				"message": "The specified parent category does not exist",
			})
		}
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

	// Update category fields
	if request.Title != "" {
		category.Title = request.Title
	}
	if request.Description != "" {
		category.Description = request.Description
	}
	if request.ParentID != nil {
		category.ParentID = request.ParentID
	}

	// Save changes
	if err := facades.Orm().Query().Save(&category); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to update category",
		})
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Parent").Where("id", category.ID).First(&category); err != nil {
		// Category was updated but failed to load relationships, still return success
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message":  "Category updated successfully",
		"category": category,
	})
}

// Destroy deletes a category
func (r *CategoryController) Destroy(ctx http.Context) http.Response {
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
			"message": "Category ID is required",
		})
	}

	// Find existing category
	var category models.Category
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&category); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Category not found",
				"message": "The requested category does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve category",
		})
	}

	// Check if category has any child categories
	childCount, err := facades.Orm().Query().Model(&models.Category{}).Where("parent_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check category children",
		})
	}
	if childCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete category",
			"message": "Category has child categories and cannot be deleted",
		})
	}

	// Check if category has any products
	productCount, err := facades.Orm().Query().Model(&models.Product{}).Where("category_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check category products",
		})
	}
	if productCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete category",
			"message": "Category has existing products and cannot be deleted",
		})
	}

	// Delete category
	if _, err := facades.Orm().Query().Delete(&category); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete category",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Category deleted successfully",
	})
}

func (r *CategoryController) Select(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	searchQuery := ctx.Request().Query("query", "")

	query := facades.Orm().Query()
	if searchQuery != "" {
		query = query.Where("title LIKE ?", "%"+searchQuery+"%")
	}

	var categories []models.Category
	if err := query.OrderBy("title", "asc").Find(&categories); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve category options",
		})
	}

	options := make([]http.Json, 0, len(categories))
	for _, c := range categories {
		options = append(options, http.Json{
			"value": c.ID,
			"label": c.Title,
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"options": options,
	})
}
