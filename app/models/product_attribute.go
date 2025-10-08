package models

import (
	"github.com/goravel/framework/database/orm"
)

type ProductAttribute struct {
	orm.Model
	ProductID  uint   `gorm:"not null;index"`
	Key        string `gorm:"size:100;not null"`
	Title      string `gorm:"size:255;not null"`
	OrderIndex int    `gorm:"not null;index"`

	// Relationships
	Product Product                 `gorm:"foreignKey:ProductID"`
	Values  []ProductAttributeValue `gorm:"foreignKey:AttributeID"`
}
