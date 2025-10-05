package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type StockRequest struct {
	orm.Model
	RequestNumber string     `gorm:"size:100;unique;not null"`
	ArticleID     uint       `gorm:"not null;index"`
	VariantID     *uint      `gorm:"index"`
	LocationID    uint       `gorm:"not null;index"`
	Quantity      float64    `gorm:"type:decimal(10,3);not null"`
	Unit          string     `gorm:"size:50"`
	Status        string     `gorm:"size:20;not null;default:'pending';index"` // pending, approved, rejected, fulfilled
	Priority      string     `gorm:"size:20;not null;default:'medium';index"`  // low, medium, high, urgent
	Notes         string     `gorm:"type:text"`
	RequestedBy   uint       `gorm:"not null;index"`
	ApprovedBy    *uint      `gorm:"index"`
	RequestedAt   time.Time  `gorm:"not null;index"`
	ApprovedAt    *time.Time

	// Relationships
	Article         Article          `gorm:"foreignKey:ArticleID"`
	Variant         *ArticleVariant  `gorm:"foreignKey:VariantID"`
	Location        StorageLocation  `gorm:"foreignKey:LocationID"`
	RequestedByUser User             `gorm:"foreignKey:RequestedBy"`
	ApprovedByUser  *User            `gorm:"foreignKey:ApprovedBy"`
}