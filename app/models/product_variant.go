package models

import (
	"github.com/goravel/framework/database/orm"
)

type ProductVariant struct {
	orm.Model
	ProductID   uint    `gorm:"not null;index"`
	Title       string  `gorm:"size:255;not null"`
	Description string  `gorm:"type:text"`
	SKU         string  `gorm:"size:100;uniqueIndex"`
	Attributes  string  `gorm:"type:json"` // JSON field for variant attributes
	PrixAchat   float64 `gorm:"type:decimal(10,2)"`
	PrixVente   float64 `gorm:"type:decimal(10,2)"`
	Unit        string  `gorm:"size:50"`
	ImageURL    string  `gorm:"size:500"`
	ImageIndex  int     `gorm:"default:0"`
	IsActive    bool    `gorm:"not null;index"`

	// Relationships
	Product           Product            `gorm:"foreignKey:ProductID"`
	OrderFabrications []OrderFabrication `gorm:"foreignKey:VariantID"`
	StockLevels       []StockLevel       `gorm:"foreignKey:VariantID"`
	StockMovements    []StockMovement    `gorm:"foreignKey:VariantID"`
}
