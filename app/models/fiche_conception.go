package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type FicheConception struct {
	orm.Model
	Reference        string `gorm:"size:100;uniqueIndex;not null;index"`
	ProductVariantID *uint  `gorm:"index"`
	Title            string `gorm:"size:255;not null;index"`
	Description      string `gorm:"type:text"`
	RequestedBy      uint   `gorm:"not null;index"`
	Status           string `gorm:"size:20;not null;default:'pending';index"` // pending, in_design, design_done, cancelled
	DesignNotes      string `gorm:"type:text"`
	ValidatedBy      *uint  `gorm:"index"`
	ValidatedAt      *time.Time

	// Relationships
	CreatedByUser   User            `gorm:"foreignKey:created_by"`
	UpdatedByUser   *User           `gorm:"foreignKey:updated_by"`
	ProductVariant  *ProductVariant `gorm:"foreignKey:product_variant_id"`
	RequestedByUser User            `gorm:"foreignKey:requested_by"`
	ValidatedByUser *User           `gorm:"foreignKey:validated_by"`
}
