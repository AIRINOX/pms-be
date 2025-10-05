package models

import (
	"github.com/goravel/framework/database/orm"
)

type Article struct {
	orm.Model
	Title         string  `gorm:"size:255;not null;index"`
	Description   string  `gorm:"type:text"`
	SKU           string  `gorm:"size:100;uniqueIndex"`
	IsRawMaterial bool    `gorm:"not null;default:false;index"`
	CategoryID    *uint   `gorm:"index"`
	LocationID    *uint   `gorm:"index"`
	PrixAchat     float64 `gorm:"type:decimal(10,2)"`
	PrixVente     float64 `gorm:"type:decimal(10,2)"`
	Unit          string  `gorm:"size:50"`
	ImageURL      string  `gorm:"size:500"`

	// Relationships
	Category         *Category          `gorm:"foreignKey:CategoryID"`
	Location         *StorageLocation   `gorm:"foreignKey:LocationID"`
	Attributes       []ArticleAttribute `gorm:"foreignKey:ArticleID"`
	Variants         []ArticleVariant   `gorm:"foreignKey:ArticleID"`
	Images           []ArticleImage     `gorm:"foreignKey:ArticleID"`
	RecipeArticles   []RecipeArticle    `gorm:"foreignKey:ArticleID"`
	OrderFabrications []OrderFabrication `gorm:"foreignKey:ArticleID"`
	StockLevels      []StockLevel       `gorm:"foreignKey:ArticleID"`
	StockMovements   []StockMovement    `gorm:"foreignKey:ArticleID"`
}