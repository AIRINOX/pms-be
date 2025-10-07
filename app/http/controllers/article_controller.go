package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type ArticleController struct {
	// Dependent services
}

func NewArticleController() *ArticleController {
	return &ArticleController{
		// Inject services
	}
}

// CreateArticleRequest represents the article creation request payload
type CreateArticleRequest struct {
	Title         string  `json:"title" form:"title" validate:"required|min_len:2|max_len:255"`
	Description   string  `json:"description" form:"description"`
	SKU           string  `json:"sku" form:"sku" validate:"max_len:100"`
	IsRawMaterial bool    `json:"is_raw_material" form:"is_raw_material"`
	CategoryID    *uint   `json:"category_id" form:"category_id"`
	LocationID    *uint   `json:"location_id" form:"location_id"`
	PrixAchat     float64 `json:"prix_achat" form:"prix_achat"`
	PrixVente     float64 `json:"prix_vente" form:"prix_vente"`
	Unit          string  `json:"unit" form:"unit" validate:"max_len:50"`
	ImageURL      string  `json:"image_url" form:"image_url" validate:"max_len:500"`
}

// UpdateArticleRequest represents the article update request payload
type UpdateArticleRequest struct {
	Title         string  `json:"title" form:"title" validate:"min_len:2|max_len:255"`
	Description   string  `json:"description" form:"description"`
	SKU           string  `json:"sku" form:"sku" validate:"max_len:100"`
	IsRawMaterial *bool   `json:"is_raw_material" form:"is_raw_material"`
	CategoryID    *uint   `json:"category_id" form:"category_id"`
	LocationID    *uint   `json:"location_id" form:"location_id"`
	PrixAchat     float64 `json:"prix_achat" form:"prix_achat"`
	PrixVente     float64 `json:"prix_vente" form:"prix_vente"`
	Unit          string  `json:"unit" form:"unit" validate:"max_len:50"`
	ImageURL      string  `json:"image_url" form:"image_url" validate:"max_len:500"`
}

// CreateAttributeRequest represents the attribute creation request
type CreateAttributeRequest struct {
	Key        string   `json:"key" form:"key" validate:"required|min_len:1|max_len:100"`
	Title      string   `json:"title" form:"title" validate:"required|min_len:2|max_len:255"`
	OrderIndex int      `json:"order_index" form:"order_index"`
	Values     []string `json:"values" form:"values"`
}

// CreateVariantRequest represents the variant creation request
type CreateVariantRequest struct {
	Title       string            `json:"title" form:"title" validate:"required|min_len:2|max_len:255"`
	Description string            `json:"description" form:"description"`
	SKU         string            `json:"sku" form:"sku" validate:"max_len:100"`
	Attributes  map[string]string `json:"attributes" form:"attributes"`
	PrixAchat   float64           `json:"prix_achat" form:"prix_achat"`
	PrixVente   float64           `json:"prix_vente" form:"prix_vente"`
	Unit        string            `json:"unit" form:"unit" validate:"max_len:50"`
	ImageURL    string            `json:"image_url" form:"image_url" validate:"max_len:500"`
	ImageIndex  int               `json:"image_index" form:"image_index"`
	IsActive    bool              `json:"is_active" form:"is_active"`
}

// CreateImageRequest represents the image upload request
type CreateImageRequest struct {
	FilePath   string `json:"file_path" form:"file_path" validate:"required|max_len:500"`
	FileName   string `json:"file_name" form:"file_name" validate:"required|max_len:255"`
	ImageIndex int    `json:"image_index" form:"image_index"`
	IsPrimary  bool   `json:"is_primary" form:"is_primary"`
}

// isMethodesOrAdmin checks if the authenticated user is methodes or admin
func (r *ArticleController) isMethodesOrAdmin(ctx http.Context) bool {
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

// Index returns a paginated list of articles with search and filtering
func (r *ArticleController) Index(ctx http.Context) http.Response {
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

	// Parse filter data
	filterTitle := ctx.Request().Query("filterData[title]", "")
	filterSKU := ctx.Request().Query("filterData[sku]", "")
	filterCategory := ctx.Request().Query("filterData[category_id]", "")
	filterLocation := ctx.Request().Query("filterData[location_id]", "")
	filterIsRawMaterial := ctx.Request().Query("filterData[is_raw_material]", "")

	query := facades.Orm().Query().With("Category", "Location")

	// Apply search filter
	if searchQuery != "" {
		query = query.Where("title LIKE ? OR description LIKE ? OR sku LIKE ?",
			"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	// Apply specific filters
	if filterTitle != "" {
		query = query.Where("title LIKE ?", "%"+filterTitle+"%")
	}
	if filterSKU != "" {
		query = query.Where("sku LIKE ?", "%"+filterSKU+"%")
	}
	if filterCategory != "" {
		query = query.Where("category_id", filterCategory)
	}
	if filterLocation != "" {
		query = query.Where("location_id", filterLocation)
	}
	if filterIsRawMaterial != "" {
		if filterIsRawMaterial == "true" {
			query = query.Where("is_raw_material", true)
		} else if filterIsRawMaterial == "false" {
			query = query.Where("is_raw_material", false)
		}
	}

	// Apply sorting
	if sortKey != "" && (sortOrder == "asc" || sortOrder == "desc") {
		query = query.OrderBy(sortKey, sortOrder)
	} else {
		query = query.OrderBy("title", "asc")
	}

	var articles []models.Article

	// Get total count
	total, err := query.Model(&models.Article{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count articles",
		})
	}

	// Get paginated results
	offset := (pageIndex - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&articles); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve articles",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"articles": articles,
		"pagination": http.Json{
			"current_page": pageIndex,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific article by ID with all relationships
func (r *ArticleController) Show(ctx http.Context) http.Response {
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
			"message": "Article ID is required",
		})
	}

	var article models.Article
	if err := facades.Orm().Query().With("Category", "Location", "Attributes.Values", "Variants", "Images").Where("id", id).First(&article); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Article not found",
				"message": "The requested article does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve article",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"article": article,
	})
}

// Store creates a new article (Step 1: Basic Info)
func (r *ArticleController) Store(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	var request CreateArticleRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate input
	validationData := map[string]any{
		"title":           request.Title,
		"description":     request.Description,
		"sku":             request.SKU,
		"is_raw_material": request.IsRawMaterial,
		"prix_achat":      request.PrixAchat,
		"prix_vente":      request.PrixVente,
		"unit":            request.Unit,
		"image_url":       request.ImageURL,
	}
	validationRules := map[string]string{
		"title":           "required|min_len:2|max_len:255",
		"description":     "max_len:1000",
		"sku":             "max_len:100",
		"is_raw_material": "bool",
		"prix_achat":      "numeric",
		"prix_vente":      "numeric",
		"unit":            "max_len:50",
		"image_url":       "max_len:500",
	}

	if request.CategoryID != nil {
		validationData["category_id"] = *request.CategoryID
		validationRules["category_id"] = "numeric"
	}

	if request.LocationID != nil {
		validationData["location_id"] = *request.LocationID
		validationRules["location_id"] = "numeric"
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

	// Check if SKU already exists (if provided)
	if request.SKU != "" {
		var existingArticle models.Article
		if err := facades.Orm().Query().Where("sku", request.SKU).FirstOrFail(&existingArticle); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "SKU already exists",
				"message": "An article with this SKU already exists",
			})
		}
	}

	// Verify category exists if provided
	if request.CategoryID != nil {
		var category models.Category
		if err := facades.Orm().Query().Where("id", *request.CategoryID).FirstOrFail(&category); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid category",
				"message": "The specified category does not exist",
			})
		}
	}

	// Verify storage location exists if provided
	if request.LocationID != nil {
		var location models.StorageLocation
		if err := facades.Orm().Query().Where("id", *request.LocationID).FirstOrFail(&location); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid storage location",
				"message": "The specified storage location does not exist",
			})
		}
	}

	// Create new article
	article := models.Article{
		Title:         request.Title,
		Description:   request.Description,
		SKU:           request.SKU,
		IsRawMaterial: request.IsRawMaterial,
		CategoryID:    request.CategoryID,
		LocationID:    request.LocationID,
		PrixAchat:     request.PrixAchat,
		PrixVente:     request.PrixVente,
		Unit:          request.Unit,
		ImageURL:      request.ImageURL,
	}

	if err := facades.Orm().Query().Create(&article); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create article",
		})
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Category", "Location").Where("id", article.ID).First(&article); err != nil {
		// Article was created but failed to load relationships, still return success
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Article created successfully",
		"article": article,
	})
}

// Update modifies an existing article
func (r *ArticleController) Update(ctx http.Context) http.Response {
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
			"message": "Article ID is required",
		})
	}

	var request UpdateArticleRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing article
	var article models.Article
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&article); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Article not found",
				"message": "The requested article does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve article",
		})
	}

	// Prepare validation rules and data
	validationData := make(map[string]any)
	validationRules := make(map[string]string)

	if request.Title != "" {
		validationData["title"] = request.Title
		validationRules["title"] = "min_len:2|max_len:255"
	}

	if request.Description != "" {
		validationData["description"] = request.Description
		validationRules["description"] = "max_len:1000"
	}

	if request.SKU != "" {
		validationData["sku"] = request.SKU
		validationRules["sku"] = "max_len:100"

		// Check if SKU already exists (excluding current article)
		var existingArticle models.Article
		if err := facades.Orm().Query().Where("sku", request.SKU).Where("id != ?", id).FirstOrFail(&existingArticle); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "SKU already exists",
				"message": "An article with this SKU already exists",
			})
		}
	}

	if request.Unit != "" {
		validationData["unit"] = request.Unit
		validationRules["unit"] = "max_len:50"
	}

	if request.ImageURL != "" {
		validationData["image_url"] = request.ImageURL
		validationRules["image_url"] = "max_len:500"
	}

	if request.CategoryID != nil {
		validationData["category_id"] = *request.CategoryID
		validationRules["category_id"] = "numeric"

		// Verify category exists
		var category models.Category
		if err := facades.Orm().Query().Where("id", *request.CategoryID).FirstOrFail(&category); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid category",
				"message": "The specified category does not exist",
			})
		}
	}

	if request.LocationID != nil {
		validationData["location_id"] = *request.LocationID
		validationRules["location_id"] = "numeric"

		// Verify storage location exists
		var location models.StorageLocation
		if err := facades.Orm().Query().Where("id", *request.LocationID).FirstOrFail(&location); err != nil {
			return ctx.Response().Status(400).Json(http.Json{
				"error":   "Invalid storage location",
				"message": "The specified storage location does not exist",
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

	// Update article fields
	if request.Title != "" {
		article.Title = request.Title
	}
	if request.Description != "" {
		article.Description = request.Description
	}
	if request.SKU != "" {
		article.SKU = request.SKU
	}
	if request.IsRawMaterial != nil {
		article.IsRawMaterial = *request.IsRawMaterial
	}
	if request.CategoryID != nil {
		article.CategoryID = request.CategoryID
	}
	if request.LocationID != nil {
		article.LocationID = request.LocationID
	}
	if request.PrixAchat != 0 {
		article.PrixAchat = request.PrixAchat
	}
	if request.PrixVente != 0 {
		article.PrixVente = request.PrixVente
	}
	if request.Unit != "" {
		article.Unit = request.Unit
	}
	if request.ImageURL != "" {
		article.ImageURL = request.ImageURL
	}

	// Save changes
	if err := facades.Orm().Query().Save(&article); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to update article",
		})
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Category", "Location").Where("id", article.ID).First(&article); err != nil {
		// Article was updated but failed to load relationships, still return success
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Article updated successfully",
		"article": article,
	})
}

// Destroy deletes an article
func (r *ArticleController) Destroy(ctx http.Context) http.Response {
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
			"message": "Article ID is required",
		})
	}

	// Find existing article
	var article models.Article
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&article); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Article not found",
				"message": "The requested article does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve article",
		})
	}

	// Check if article has any orders
	orderCount, err := facades.Orm().Query().Model(&models.OrderFabrication{}).Where("article_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check article orders",
		})
	}
	if orderCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete article",
			"message": "Article has existing fabrication orders and cannot be deleted",
		})
	}

	// Check if article has any stock levels
	stockLevelCount, err := facades.Orm().Query().Model(&models.StockLevel{}).Where("article_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check article stock levels",
		})
	}
	if stockLevelCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete article",
			"message": "Article has existing stock levels and cannot be deleted",
		})
	}

	// Delete related data first
	// Delete attributes and their values
	var attributes []models.ArticleAttribute
	if err := facades.Orm().Query().Where("article_id", id).Find(&attributes); err == nil {
		for _, attr := range attributes {
			// Delete attribute values
			facades.Orm().Query().Where("attribute_id", attr.ID).Delete(&models.ArticleAttributeValue{})
		}
		// Delete attributes
		facades.Orm().Query().Where("article_id", id).Delete(&models.ArticleAttribute{})
	}

	// Delete variants
	facades.Orm().Query().Where("article_id", id).Delete(&models.ArticleVariant{})

	// Delete images
	facades.Orm().Query().Where("article_id", id).Delete(&models.ArticleImage{})

	// Delete recipes
	facades.Orm().Query().Where("article_id", id).Delete(&models.RecipeArticle{})

	// Delete article
	if _, err := facades.Orm().Query().Delete(&article); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete article",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Article deleted successfully",
	})
}

// CreateAttribute creates attributes for an article (Step 2: Attributes definition)
func (r *ArticleController) CreateAttribute(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var request CreateAttributeRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate input
	validator, err := facades.Validation().Make(map[string]any{
		"key":         request.Key,
		"title":       request.Title,
		"order_index": request.OrderIndex,
	}, map[string]string{
		"key":         "required|min_len:1|max_len:100",
		"title":       "required|min_len:2|max_len:255",
		"order_index": "numeric",
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

	// Verify article exists
	var article models.Article
	if err := facades.Orm().Query().Where("id", articleID).FirstOrFail(&article); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Article not found",
			"message": "The specified article does not exist",
		})
	}

	// Check if attribute key already exists for this article
	var existingAttribute models.ArticleAttribute
	if err := facades.Orm().Query().Where("article_id", articleID).Where("key", request.Key).FirstOrFail(&existingAttribute); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Attribute key already exists",
			"message": "An attribute with this key already exists for this article",
		})
	}

	// Convert articleID to uint
	articleIDUint, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid article ID",
			"message": "Article ID must be a valid number",
		})
	}

	// Create new attribute
	attribute := models.ArticleAttribute{
		ArticleID:  uint(articleIDUint),
		Key:        request.Key,
		Title:      request.Title,
		OrderIndex: request.OrderIndex,
	}

	if err := facades.Orm().Query().Create(&attribute); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create attribute",
		})
	}

	// Create attribute values if provided
	if len(request.Values) > 0 {
		for i, value := range request.Values {
			attributeValue := models.ArticleAttributeValue{
				AttributeID: attribute.ID,
				Value:       value,
				OrderIndex:  i,
				IsActive:    true,
			}
			facades.Orm().Query().Create(&attributeValue)
		}
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Values").Where("id", attribute.ID).First(&attribute); err != nil {
		// Attribute was created but failed to load relationships, still return success
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message":   "Attribute created successfully",
		"attribute": attribute,
	})
}

// CreateImages uploads multiple images for an article (Step 3: Upload multiple images)
func (r *ArticleController) CreateImages(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var request []CreateImageRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Verify article exists
	var article models.Article
	if err := facades.Orm().Query().Where("id", articleID).FirstOrFail(&article); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Article not found",
			"message": "The specified article does not exist",
		})
	}

	// Convert articleID to uint
	articleIDUint, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid article ID",
			"message": "Article ID must be a valid number",
		})
	}

	var createdImages []models.ArticleImage

	// Create images
	for _, imageReq := range request {
		// Validate each image request
		validator, err := facades.Validation().Make(map[string]any{
			"file_path":   imageReq.FilePath,
			"file_name":   imageReq.FileName,
			"image_index": imageReq.ImageIndex,
			"is_primary":  imageReq.IsPrimary,
		}, map[string]string{
			"file_path":   "required|max_len:500",
			"file_name":   "required|max_len:255",
			"image_index": "numeric",
			"is_primary":  "bool",
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

		// If this image is set as primary, unset other primary images
		if imageReq.IsPrimary {
			facades.Orm().Query().Model(&models.ArticleImage{}).Where("article_id", articleID).Update("is_primary", false)
		}

		image := models.ArticleImage{
			ArticleID:  uint(articleIDUint),
			FilePath:   imageReq.FilePath,
			FileName:   imageReq.FileName,
			ImageIndex: imageReq.ImageIndex,
			IsPrimary:  imageReq.IsPrimary,
		}

		if err := facades.Orm().Query().Create(&image); err != nil {
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Database error",
				"message": "Failed to create image",
			})
		}

		createdImages = append(createdImages, image)
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Images created successfully",
		"images":  createdImages,
	})
}

// CreateVariant creates a variant for an article (Step 4: Add product variants)
func (r *ArticleController) CreateVariant(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var request CreateVariantRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Validate input
	validator, err := facades.Validation().Make(map[string]any{
		"title":       request.Title,
		"description": request.Description,
		"sku":         request.SKU,
		"prix_achat":  request.PrixAchat,
		"prix_vente":  request.PrixVente,
		"unit":        request.Unit,
		"image_url":   request.ImageURL,
		"image_index": request.ImageIndex,
		"is_active":   request.IsActive,
	}, map[string]string{
		"title":       "required|min_len:2|max_len:255",
		"description": "max_len:1000",
		"sku":         "max_len:100",
		"prix_achat":  "numeric",
		"prix_vente":  "numeric",
		"unit":        "max_len:50",
		"image_url":   "max_len:500",
		"image_index": "numeric",
		"is_active":   "bool",
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

	// Verify article exists
	var article models.Article
	if err := facades.Orm().Query().Where("id", articleID).FirstOrFail(&article); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Article not found",
			"message": "The specified article does not exist",
		})
	}

	// Check if SKU already exists (if provided)
	if request.SKU != "" {
		var existingVariant models.ArticleVariant
		if err := facades.Orm().Query().Where("sku", request.SKU).FirstOrFail(&existingVariant); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "SKU already exists",
				"message": "A variant with this SKU already exists",
			})
		}
	}

	// Convert articleID to uint
	articleIDUint, err := strconv.ParseUint(articleID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid article ID",
			"message": "Article ID must be a valid number",
		})
	}

	// Convert attributes map to JSON
	attributesJSON, err := json.Marshal(request.Attributes)
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "JSON encoding error",
			"message": "Failed to encode attributes",
		})
	}

	// Create new variant
	variant := models.ArticleVariant{
		ArticleID:   uint(articleIDUint),
		Title:       request.Title,
		Description: request.Description,
		SKU:         request.SKU,
		Attributes:  string(attributesJSON),
		PrixAchat:   request.PrixAchat,
		PrixVente:   request.PrixVente,
		Unit:        request.Unit,
		ImageURL:    request.ImageURL,
		ImageIndex:  request.ImageIndex,
		IsActive:    request.IsActive,
	}

	if err := facades.Orm().Query().Create(&variant); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create variant",
		})
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Variant created successfully",
		"variant": variant,
	})
}

// GetAttributes returns all attributes for an article
func (r *ArticleController) GetAttributes(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var attributes []models.ArticleAttribute
	if err := facades.Orm().Query().With("Values").Where("article_id", articleID).OrderBy("order_index").Find(&attributes); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve attributes",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"attributes": attributes,
	})
}

// GetVariants returns all variants for an article
func (r *ArticleController) GetVariants(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var variants []models.ArticleVariant
	if err := facades.Orm().Query().Where("article_id", articleID).OrderBy("title").Find(&variants); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve variants",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"variants": variants,
	})
}

// GetImages returns all images for an article
func (r *ArticleController) GetImages(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	articleID := ctx.Request().Route("id")
	if articleID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Article ID is required",
		})
	}

	var images []models.ArticleImage
	if err := facades.Orm().Query().Where("article_id", articleID).OrderBy("image_index").Find(&images); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve images",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"images": images,
	})
}
