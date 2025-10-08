package models

import (
	"github.com/goravel/framework/database/orm"
)

type StorageLocation struct {
	orm.Model
	Name        string `gorm:"size:255;not null;index"`
	Description string `gorm:"type:text"`

	// Relationships
	Products       []Product       `gorm:"foreignKey:LocationID"`
	StockLevels    []StockLevel    `gorm:"foreignKey:LocationID"`
	StockMovements []StockMovement `gorm:"foreignKey:LocationID"`
}
