package models

import (
	"github.com/goravel/framework/database/orm"
)

type ClientSite struct {
	orm.Model
	ClientID     uint   `gorm:"not null;index"`
	Title        string `gorm:"size:255;not null;index"`
	Address      string `gorm:"type:text"`
	ContactName  string `gorm:"size:255"`
	ContactPhone string `gorm:"size:20"`

	// Relationships
	Client            Client            `gorm:"foreignKey:ClientID"`
	OrderFabrications []OrderFabrication `gorm:"foreignKey:ClientSiteID"`
}