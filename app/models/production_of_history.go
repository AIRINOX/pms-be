package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type ProductionOfHistory struct {
	orm.Model
	OrderFabricationID uint      `gorm:"not null;index"`
	OperationID        uint      `gorm:"not null;index"`
	UserID             uint      `gorm:"not null;index"`
	Status             string    `gorm:"size:50;not null;index"`
	Notes              string    `gorm:"type:text"`
	StatusAt           time.Time `gorm:"not null;index"`

	// Relationships
	OrderFabrication OrderFabrication `gorm:"foreignKey:OrderFabricationID"`
	Operation        Operation        `gorm:"foreignKey:OperationID"`
	User             User             `gorm:"foreignKey:UserID"`
}