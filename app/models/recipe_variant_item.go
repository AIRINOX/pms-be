package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeVariantItem struct {
	orm.Model
	RecipeVariantID   uint    `gorm:"not null;index"`
	MaterialVariantID uint    `gorm:"not null;index"`
	Quantity          float64 `gorm:"type:decimal(10,3);not null"`
	Notes             string  `gorm:"type:text"`

	// Relationships
	RecipeVariant   RecipeVariant  `gorm:"foreignKey:RecipeVariantID"`
	MaterialVariant ArticleVariant `gorm:"foreignKey:MaterialVariantID"`
}
