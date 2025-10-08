package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeVariant struct {
	orm.Model
	ProductID      uint    `gorm:"not null;index"`
	VariantID      uint    `gorm:"not null;index"`
	OutputQuantity float64 `gorm:"type:decimal(10,3);not null"`
	Notes          string  `gorm:"type:text"`

	// Relationships
	Product            Product             `gorm:"foreignKey:ProductID"`
	Variant            ProductVariant      `gorm:"foreignKey:VariantID"`
	RecipeVariantItems []RecipeVariantItem `gorm:"foreignKey:RecipeVariantID"`
}
