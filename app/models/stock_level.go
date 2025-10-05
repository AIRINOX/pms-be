package models

import (
	"github.com/goravel/framework/database/orm"
)

type StockLevel struct {
	orm.Model
	ArticleID  uint    `gorm:"not null;index"`
	VariantID  *uint   `gorm:"index"`
	LocationID uint    `gorm:"not null;index"`
	Quantity   float64 `gorm:"not null;index"`
	Unit       string  `gorm:"size:50"`

	// Relationships
	Article  Article         `gorm:"foreignKey:ArticleID"`
	Variant  *ArticleVariant `gorm:"foreignKey:VariantID"`
	Location StorageLocation `gorm:"foreignKey:LocationID"`
}
