package models

import (
	"github.com/goravel/framework/database/orm"
)

type ArticleAttribute struct {
	orm.Model
	ArticleID  uint   `gorm:"not null;index"`
	Key        string `gorm:"size:100;not null"`
	Title      string `gorm:"size:255;not null"`
	OrderIndex int    `gorm:"not null;index"`

	// Relationships
	Article Article                 `gorm:"foreignKey:ArticleID"`
	Values  []ArticleAttributeValue `gorm:"foreignKey:AttributeID"`
}
