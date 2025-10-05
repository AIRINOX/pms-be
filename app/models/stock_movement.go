package models

import (
	"github.com/goravel/framework/database/orm"
)

type StockMovement struct {
	orm.Model
	ArticleID     uint    `gorm:"not null;index"`
	VariantID     *uint   `gorm:"index"`
	LocationID    uint    `gorm:"not null;index"`
	MovementType  string  `gorm:"size:20;not null;index"` // in, out, adjustment
	Quantity      float64 `gorm:"not null"`
	Unit          string  `gorm:"size:50"`
	ReferenceType string  `gorm:"size:50;index"`
	ReferenceID   *uint   `gorm:"index"`
	Notes         string  `gorm:"type:text"`
	CreatedBy     uint    `gorm:"not null;index"`

	// Relationships
	Article  Article         `gorm:"foreignKey:ArticleID"`
	Variant  *ArticleVariant `gorm:"foreignKey:VariantID"`
	Location StorageLocation `gorm:"foreignKey:LocationID"`
	Creator  User            `gorm:"foreignKey:CreatedBy"`
}
