package models

import (
	"github.com/goravel/framework/database/orm"
)

type TechnicalDocument struct {
	orm.Model
	Title       string `gorm:"size:255;not null;index"`
	Description string `gorm:"type:text"`
	FilePath    string `gorm:"size:500;not null"`
	FileName    string `gorm:"size:255;not null"`
	FileType    string `gorm:"size:50;not null;index"`
	FileSize    uint64 `gorm:"not null"`
	ArticleID   *uint  `gorm:"index"`
	VariantID   *uint  `gorm:"index"`
	UploadedBy  uint   `gorm:"not null;index"`

	// Relationships
	Article        *Article        `gorm:"foreignKey:ArticleID"`
	Variant        *ArticleVariant `gorm:"foreignKey:VariantID"`
	UploadedByUser User            `gorm:"foreignKey:UploadedBy"`
}