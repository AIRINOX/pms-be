package models

import (
	"github.com/goravel/framework/database/orm"
)

type ProductImage struct {
	orm.Model
	ProductID  uint   `gorm:"not null;index"`
	FileUrl    string `gorm:"size:500;not null"`
	FileName   string `gorm:"size:255;not null"`
	ImageIndex int    `gorm:"not null;index"`
	IsPrimary  bool   `gorm:"not null;default:false;index"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
