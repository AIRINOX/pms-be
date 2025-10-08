package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type OrderFabrication struct {
	orm.Model
	OrderNumber  string     `gorm:"size:100;uniqueIndex;not null"`
	ProductID    uint       `gorm:"not null;index"`
	VariantID    *uint      `gorm:"index"`
	Quantity     float64    `gorm:"not null"`
	ClientID     uint       `gorm:"not null;index"`
	ClientSiteID *uint      `gorm:"index"`
	Status       string     `gorm:"size:50;not null;default:'pending';index"` // pending, in_progress, completed, cancelled, on_hold
	Priority     string     `gorm:"size:20;not null;default:'normal';index"`  // low, normal, high, urgent
	DeadlineDate *time.Time `gorm:"index"`
	Notes        string     `gorm:"type:text"`
	CreatedBy    uint       `gorm:"not null;index"`

	// Relationships
	Product    Product         `gorm:"foreignKey:ProductID"`
	Variant    *ProductVariant `gorm:"foreignKey:VariantID"`
	Client     Client          `gorm:"foreignKey:ClientID"`
	ClientSite *ClientSite     `gorm:"foreignKey:ClientSiteID"`
	Creator    User            `gorm:"foreignKey:CreatedBy"`
}
