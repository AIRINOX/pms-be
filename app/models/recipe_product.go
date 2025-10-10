package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeProduct struct {
	orm.Model
	ProductID         uint   `gorm:"not null;index"`
	MaterialProductID uint   `gorm:"not null;index"`
	Notes             string `gorm:"type:text"`

	// Relationships
	Product         Product `gorm:"foreignKey:product_id"`
	MaterialProduct Product `gorm:"foreignKey:material_product_id"`
}
