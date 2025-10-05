package models

import (
	"github.com/goravel/framework/database/orm"
)

type ProductionMaterialRequirement struct {
	orm.Model
	OrderFabricationID uint    `gorm:"not null;index"`
	MaterialVariantID  uint    `gorm:"not null;index"`
	RequiredQuantity   float64 `gorm:"type:decimal(10,3);not null"`
	StockQuantity      float64 `gorm:"type:decimal(10,3);default:0"`
	RequestQuantity    float64 `gorm:"type:decimal(10,3);default:0"`
	Unit               string  `gorm:"size:50"`
	Status             string  `gorm:"size:20;not null;default:'pending';index"` // pending, requested, available, consumed

	// Relationships
	OrderFabrication OrderFabrication `gorm:"foreignKey:OrderFabricationID"`
	MaterialVariant  ArticleVariant   `gorm:"foreignKey:MaterialVariantID"`
}