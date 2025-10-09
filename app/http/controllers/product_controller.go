package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"

	"pms/app/models"
)

type ProductController struct {
	// Dependent services
}

func NewProductController() *ProductController {
	return &ProductController{
		// Inject services
	}
}

// CreateProductRequest represents the product creation request payload
type CreateProductRequest struct {
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

// UpdateProductRequest represents the product update request payload
type UpdateProductRequest struct {
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

// ProductAttributeOption represents an attribute and its values
type ProductAttributeOption struct {
	AttributeID     uint     `json:"attribute_id"`
	AttributeName   string   `json:"attribute_name"`
	AttributeValues []string `json:"attribute_values"`
}

// ProductVariantOption represents a variant option
type ProductVariantOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductVariantRequest represents a product variant in the request
type ProductVariantRequest struct {
	SKU      string                 `json:"sku"`
	Price    float64                `json:"price"`
	Quantity int                    `json:"quantity"`
	Options  []ProductVariantOption `json:"options"`
}

// ProductImageRequest represents a product image in the request
type ProductImageRequest struct {
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

// CompleteProductUpdateRequest represents the full product update request payload
type CompleteProductUpdateRequest struct {
	Title         string                   `json:"title" validate:"required|min_len:2|max_len:255"`
	CategoryID    uint                     `json:"category_id"`
	LocationID    uint                     `json:"location_id"`
	PrixAchat     float64                  `json:"prix_achat"`
	PrixVente     float64                  `json:"prix_vente"`
	IsRawMaterial bool                     `json:"is_raw_material"`
	Attributes    []ProductAttributeOption `json:"attributes"`
	Variants      []ProductVariantRequest  `json:"variants"`
	Images        []ProductImageRequest    `json:"images"`
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
func (r *ProductController) isMethodesOrAdmin(ctx http.Context) bool {
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

// Index returns a paginated list of products with search and filtering
func (r *ProductController) Index(ctx http.Context) http.Response {
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

	query := facades.Orm().Query().With("Category").With("Location")

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

	var products []models.Product

	// Get total count
	total, err := query.Model(&models.Product{}).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to count products",
		})
	}

	// Get paginated results
	offset := (pageIndex - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&products); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve products",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"products": products,
		"pagination": http.Json{
			"current_page": pageIndex,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a specific product by ID with all relationships
func (r *ProductController) Show(ctx http.Context) http.Response {
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
			"message": "Product ID is required",
		})
	}

	var product models.Product
	if err := facades.Orm().Query().With("Category").With("Location").With("Attributes.Values").With("Variants").With("Images").Where("id", id).First(&product); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve product",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"product": product,
	})
}

// Store creates a new product (Step 1: Basic Info)
func (r *ProductController) Store(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	var request CreateProductRequest

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
		var existingProduct models.Product
		if err := facades.Orm().Query().Where("sku", request.SKU).FirstOrFail(&existingProduct); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "SKU already exists",
				"message": "An product with this SKU already exists",
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

	// Create new product
	product := models.Product{
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

	if err := facades.Orm().Query().Create(&product); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to create product",
		})
	}

	// Load relationships for response
	if err := facades.Orm().Query().With("Category").With("Location").Where("id", product.ID).First(&product); err != nil {
		// Product was created but failed to load relationships, still return success
	}

	return ctx.Response().Status(201).Json(http.Json{
		"message": "Product created successfully",
		"product": product,
	})
}

// Update modifies an existing product with all its attributes, variants and images
func (r *ProductController) Update(ctx http.Context) http.Response {
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
			"message": "Product ID is required",
		})
	}

	var request CompleteProductUpdateRequest

	// Validate request
	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request data",
			"message": err.Error(),
		})
	}

	// Find existing product
	var product models.Product
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&product); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve product",
		})
	}

	// Verify category exists
	var category models.Category
	if err := facades.Orm().Query().Where("id", request.CategoryID).FirstOrFail(&category); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid category",
			"message": "The specified category does not exist",
		})
	}

	// Verify storage location exists
	var location models.StorageLocation
	if err := facades.Orm().Query().Where("id", request.LocationID).FirstOrFail(&location); err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid storage location",
			"message": "The specified storage location does not exist",
		})
	}

	// Start a transaction
	tx, err := facades.Orm().Query().Begin()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to start transaction",
		})
	}

	// Update product basic info
	product.Title = request.Title
	product.CategoryID = &request.CategoryID
	product.LocationID = &request.LocationID
	product.PrixAchat = request.PrixAchat
	product.PrixVente = request.PrixVente
	product.IsRawMaterial = request.IsRawMaterial

	if err := tx.Save(&product); err != nil {
		tx.Rollback()
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to update product",
		})
	}

	// Delete existing attributes and their values
	var existingAttributes []models.ProductAttribute
	if err := tx.Where("product_id", id).Find(&existingAttributes); err == nil {
		for _, attr := range existingAttributes {
			// Delete attribute values
			if _, err := tx.Where("attribute_id", attr.ID).Delete(&models.ProductAttributeValue{}); err != nil {
				tx.Rollback()
				return ctx.Response().Status(500).Json(http.Json{
					"error":   "Database error",
					"message": "Failed to delete attribute values",
				})
			}
		}
		// Delete attributes
		if _, err := tx.Where("product_id", id).Delete(&models.ProductAttribute{}); err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Database error",
				"message": "Failed to delete attributes",
			})
		}
	}

	// Create new attributes and their values
	for i, attr := range request.Attributes {
		// Create attribute
		attribute := models.ProductAttribute{
			ProductID:  product.ID,
			Key:        attr.AttributeName,
			Title:      attr.AttributeName,
			OrderIndex: i,
		}

		if err := tx.Create(&attribute); err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Database error",
				"message": "Failed to create attribute",
			})
		}

		// Create attribute values
		for j, value := range attr.AttributeValues {
			if value == "" {
				continue
			}

			attributeValue := models.ProductAttributeValue{
				AttributeID: attribute.ID,
				Value:       value,
				OrderIndex:  j,
				IsActive:    true,
			}

			if err := tx.Create(&attributeValue); err != nil {
				tx.Rollback()
				return ctx.Response().Status(500).Json(http.Json{
					"error":   "Database error",
					"message": "Failed to create attribute value",
				})
			}
		}
	}

	// Delete existing variants
	if _, err := tx.Where("product_id", id).Delete(&models.ProductVariant{}); err != nil {
		tx.Rollback()
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete variants",
		})
	}

	// Create new variants
	for _, variant := range request.Variants {
		// Convert options to JSON
		optionsMap := make(map[string]string)
		for _, option := range variant.Options {
			optionsMap[option.Name] = option.Value
		}

		optionsJSON, err := json.Marshal(optionsMap)
		if err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "JSON error",
				"message": "Failed to encode variant options",
			})
		}

		// Create variant
		newVariant := models.ProductVariant{
			ProductID:  product.ID,
			Title:      product.Title,
			SKU:        variant.SKU,
			Attributes: string(optionsJSON),
			PrixVente:  variant.Price,
			IsActive:   true,
		}

		if err := tx.Create(&newVariant); err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Database error",
				"message": "Failed to create variant",
			})
		}
	}

	// Delete existing images
	if _, err := tx.Where("product_id", id).Delete(&models.ProductImage{}); err != nil {
		tx.Rollback()
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete images",
		})
	}

	// Create new images
	for i, image := range request.Images {
		// Skip empty images
		if image.URL == "" || image.URL == "data:image/jpeg;base64," {
			continue
		}
		// i want save url as file theen will put the url of file in file path
		// i will use file name to save file in storage
		fileName := "image_" + strconv.Itoa(i+1) + ".jpg"
		filePath := "storage/products/" + fileName
		// i will save file in storage
		if err := r.saveBase64Image(image.URL, filePath); err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "File error",
				"message": "Failed to save image",
			})
		}

		newImage := models.ProductImage{
			ProductID:  product.ID,
			FilePath:   filePath,
			FileName:   fileName,
			ImageIndex: i,
			IsPrimary:  image.IsPrimary,
		}

		if err := tx.Create(&newImage); err != nil {
			tx.Rollback()
			return ctx.Response().Status(500).Json(http.Json{
				"error":   "Database error",
				"message": "Failed to create image",
			})
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to commit transaction",
		})
	}

	// Load the updated product with all relationships for response
	var updatedProduct models.Product
	if err := facades.Orm().Query().With("Category").With("Location").With("Attributes.Values").With("Variants").With("Images").Where("id", id).First(&updatedProduct); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve updated product",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Product updated successfully",
		"product": updatedProduct,
	})
}

// Destroy deletes an product
func (r *ProductController) Destroy(ctx http.Context) http.Response {
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
			"message": "Product ID is required",
		})
	}

	// Find existing product
	var product models.Product
	if err := facades.Orm().Query().Where("id", id).FirstOrFail(&product); err != nil {
		if errors.Is(err, errors.OrmRecordNotFound) {
			return ctx.Response().Status(404).Json(http.Json{
				"error":   "Product not found",
				"message": "The requested product does not exist",
			})
		}
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve product",
		})
	}

	// Check if product has any orders
	orderCount, err := facades.Orm().Query().Model(&models.OrderFabrication{}).Where("product_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check product orders",
		})
	}
	if orderCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete product",
			"message": "Product has existing fabrication orders and cannot be deleted",
		})
	}

	// Check if product has any stock levels
	stockLevelCount, err := facades.Orm().Query().Model(&models.StockLevel{}).Where("product_id", id).Count()
	if err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to check product stock levels",
		})
	}
	if stockLevelCount > 0 {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Cannot delete product",
			"message": "Product has existing stock levels and cannot be deleted",
		})
	}

	// Delete related data first
	// Delete attributes and their values
	var attributes []models.ProductAttribute
	if err := facades.Orm().Query().Where("product_id", id).Find(&attributes); err == nil {
		for _, attr := range attributes {
			// Delete attribute values
			facades.Orm().Query().Where("attribute_id", attr.ID).Delete(&models.ProductAttributeValue{})
		}
		// Delete attributes
		facades.Orm().Query().Where("product_id", id).Delete(&models.ProductAttribute{})
	}

	// Delete variants
	facades.Orm().Query().Where("product_id", id).Delete(&models.ProductVariant{})

	// Delete images
	facades.Orm().Query().Where("product_id", id).Delete(&models.ProductImage{})

	// Delete recipes
	facades.Orm().Query().Where("product_id", id).Delete(&models.RecipeProduct{})

	// Delete product
	if _, err := facades.Orm().Query().Delete(&product); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to delete product",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"message": "Product deleted successfully",
	})
}

// CreateAttribute creates attributes for an product (Step 2: Attributes definition)
func (r *ProductController) CreateAttribute(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
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

	// Verify product exists
	var product models.Product
	if err := facades.Orm().Query().Where("id", productID).FirstOrFail(&product); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Product not found",
			"message": "The specified product does not exist",
		})
	}

	// Check if attribute key already exists for this product
	var existingAttribute models.ProductAttribute
	if err := facades.Orm().Query().Where("product_id", productID).Where("key", request.Key).FirstOrFail(&existingAttribute); err == nil {
		return ctx.Response().Status(409).Json(http.Json{
			"error":   "Attribute key already exists",
			"message": "An attribute with this key already exists for this product",
		})
	}

	// Convert productID to uint
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid number",
		})
	}

	// Create new attribute
	attribute := models.ProductAttribute{
		ProductID:  uint(productIDUint),
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
			attributeValue := models.ProductAttributeValue{
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

// saveBase64Image decodes a base64 image string and saves it to the specified file path
func (r *ProductController) saveBase64Image(base64String string, filePath string) error {
	// Extract the base64 data from the data URL
	parts := strings.Split(base64String, ",")
	if len(parts) != 2 {
		return errors.New("invalid base64 image format")
	}

	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the file
	return ioutil.WriteFile(filePath, data, 0644)
}

// CreateImages uploads multiple images for an product (Step 3: Upload multiple images)
func (r *ProductController) CreateImages(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
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

	// Verify product exists
	var product models.Product
	if err := facades.Orm().Query().Where("id", productID).FirstOrFail(&product); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Product not found",
			"message": "The specified product does not exist",
		})
	}

	// Convert productID to uint
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid number",
		})
	}

	var createdImages []models.ProductImage

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
			facades.Orm().Query().Model(&models.ProductImage{}).Where("product_id", productID).Update("is_primary", false)
		}

		image := models.ProductImage{
			ProductID:  uint(productIDUint),
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

// CreateVariant creates a variant for an product (Step 4: Add product variants)
func (r *ProductController) CreateVariant(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
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

	// Verify product exists
	var product models.Product
	if err := facades.Orm().Query().Where("id", productID).FirstOrFail(&product); err != nil {
		return ctx.Response().Status(404).Json(http.Json{
			"error":   "Product not found",
			"message": "The specified product does not exist",
		})
	}

	// Check if SKU already exists (if provided)
	if request.SKU != "" {
		var existingVariant models.ProductVariant
		if err := facades.Orm().Query().Where("sku", request.SKU).FirstOrFail(&existingVariant); err == nil {
			return ctx.Response().Status(409).Json(http.Json{
				"error":   "SKU already exists",
				"message": "A variant with this SKU already exists",
			})
		}
	}

	// Convert productID to uint
	productIDUint, err := strconv.ParseUint(productID, 10, 32)
	if err != nil {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid product ID",
			"message": "Product ID must be a valid number",
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
	variant := models.ProductVariant{
		ProductID:   uint(productIDUint),
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

// GetAttributes returns all attributes for an product
func (r *ProductController) GetAttributes(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
		})
	}

	var attributes []models.ProductAttribute
	if err := facades.Orm().Query().With("Values").Where("product_id", productID).OrderBy("order_index").Find(&attributes); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve attributes",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"attributes": attributes,
	})
}

// GetVariants returns all variants for an product
func (r *ProductController) GetVariants(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
		})
	}

	var variants []models.ProductVariant
	if err := facades.Orm().Query().Where("product_id", productID).OrderBy("title").Find(&variants); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve variants",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"variants": variants,
	})
}

// GetImages returns all images for an product
func (r *ProductController) GetImages(ctx http.Context) http.Response {
	if !r.isMethodesOrAdmin(ctx) {
		return ctx.Response().Status(403).Json(http.Json{
			"error":   "Forbidden",
			"message": "Methodes or Admin access required",
		})
	}

	productID := ctx.Request().Route("id")
	if productID == "" {
		return ctx.Response().Status(400).Json(http.Json{
			"error":   "Invalid request",
			"message": "Product ID is required",
		})
	}

	var images []models.ProductImage
	if err := facades.Orm().Query().Where("product_id", productID).OrderBy("image_index").Find(&images); err != nil {
		return ctx.Response().Status(500).Json(http.Json{
			"error":   "Database error",
			"message": "Failed to retrieve images",
		})
	}

	return ctx.Response().Status(200).Json(http.Json{
		"images": images,
	})
}
