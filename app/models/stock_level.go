package models

import (
	"github.com/goravel/framework/database/orm"
)

type StockLevel struct {
	orm.Model
	ProductID  uint    `gorm:"not null;index"`
	VariantID  *uint   `gorm:"index"`
	LocationID uint    `gorm:"not null;index"`
	Quantity   float64 `gorm:"not null;index"`
	Unit       string  `gorm:"size:50"`

	// Relationships
	Product  Product         `gorm:"foreignKey:ProductID"`
	Variant  *ProductVariant `gorm:"foreignKey:VariantID"`
	Location StorageLocation `gorm:"foreignKey:LocationID"`
}
