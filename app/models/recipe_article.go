package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeProduct struct {
	orm.Model
	Title       string `gorm:"size:255;not null;index"`
	Description string `gorm:"type:text"`
	ProductID   uint   `gorm:"not null;index"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
