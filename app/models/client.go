package models

import (
	"github.com/goravel/framework/database/orm"
)

type Client struct {
	orm.Model
	Name    string `gorm:"size:255;not null;index"`
	Phone   string `gorm:"size:20"`
	Email   string `gorm:"size:255;index"`
	Address string `gorm:"type:text"`

	// Relationships
	ClientSites []ClientSite `gorm:"foreignKey:ClientID"`
}
