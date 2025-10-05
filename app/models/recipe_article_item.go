package models

import (
	"github.com/goravel/framework/database/orm"
)

type RecipeArticleItem struct {
	orm.Model
	RecipeArticleID   uint   `gorm:"not null;index"`
	MaterialArticleID uint   `gorm:"not null;index"`
	Notes             string `gorm:"type:text"`

	// Relationships
	RecipeArticle   RecipeArticle `gorm:"foreignKey:RecipeArticleID"`
	MaterialArticle Article       `gorm:"foreignKey:MaterialArticleID"`
}
