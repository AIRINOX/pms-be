package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeProductItem struct {
	orm.Model
	RecipeProductID   uint   `gorm:"not null;index"`
	MaterialProductID uint   `gorm:"not null;index"`
	Notes             string `gorm:"type:text"`

	// Relationships
	RecipeProduct   RecipeProduct `gorm:"foreignKey:RecipeProductID"`
	MaterialProduct Product       `gorm:"foreignKey:MaterialProductID"`
}
