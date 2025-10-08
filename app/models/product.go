package models

import (
	"github.com/goravel/framework/database/orm"
)

type Product struct {
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
	Category          *Category          `gorm:"foreignKey:CategoryID"`
	Location          *StorageLocation   `gorm:"foreignKey:LocationID"`
	Attributes        []ProductAttribute `gorm:"foreignKey:ProductID"`
	Variants          []ProductVariant   `gorm:"foreignKey:ProductID"`
	Images            []ProductImage     `gorm:"foreignKey:ProductID"`
	RecipeProducts    []RecipeProduct    `gorm:"foreignKey:ProductID"`
	OrderFabrications []OrderFabrication `gorm:"foreignKey:ProductID"`
	StockLevels       []StockLevel       `gorm:"foreignKey:ProductID"`
	StockMovements    []StockMovement    `gorm:"foreignKey:ProductID"`
}
