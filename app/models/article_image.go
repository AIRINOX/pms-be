package models

import (
	"github.com/goravel/framework/database/orm"
)

type ArticleImage struct {
	orm.Model
	ArticleID  uint   `gorm:"not null;index"`
	FilePath   string `gorm:"size:500;not null"`
	FileName   string `gorm:"size:255;not null"`
	ImageIndex int    `gorm:"not null;index"`
	IsPrimary  bool   `gorm:"not null;default:false;index"`

	// Relationships
	Article Article `gorm:"foreignKey:ArticleID"`
}
