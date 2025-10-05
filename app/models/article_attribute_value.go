package models

import (
	"github.com/goravel/framework/database/orm"
)

type ArticleAttributeValue struct {
	orm.Model
	AttributeID uint   `gorm:"not null;index"`
	Value       string `gorm:"size:255;not null"`
	OrderIndex  int    `gorm:"not null;index"`
	IsActive    bool   `gorm:"not null;default:true;index"`

	// Relationships
	Attribute ArticleAttribute `gorm:"foreignKey:AttributeID"`
}
