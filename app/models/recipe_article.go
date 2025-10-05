package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeArticle struct {
	orm.Model
	Title       string `gorm:"size:255;not null;index"`
	Description string `gorm:"type:text"`
	ArticleID   uint   `gorm:"not null;index"`

	// Relationships
	Article Article `gorm:"foreignKey:ArticleID"`
}