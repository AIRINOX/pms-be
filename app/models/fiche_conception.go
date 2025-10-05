package models

import (
	"github.com/goravel/framework/database/orm"
)

type FicheConception struct {
	orm.Model
	Title         string `gorm:"size:255;not null;index"`
	Description   string `gorm:"type:text"`
	ArticleID     uint   `gorm:"not null;index"`
	VariantID     *uint  `gorm:"index"`
	Specifications string `gorm:"type:json"`
	Materials     string `gorm:"type:json"`
	Dimensions    string `gorm:"type:json"`
	Notes         string `gorm:"type:text"`
	CreatedBy     uint   `gorm:"not null;index"`
	UpdatedBy     *uint  `gorm:"index"`

	// Relationships
	Article       Article         `gorm:"foreignKey:ArticleID"`
	Variant       *ArticleVariant `gorm:"foreignKey:VariantID"`
	CreatedByUser User            `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User           `gorm:"foreignKey:UpdatedBy"`
}