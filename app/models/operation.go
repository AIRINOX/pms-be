package models

import (
	"github.com/goravel/framework/database/orm"
)

type Operation struct {
	orm.Model
	Key        string `gorm:"uniqueIndex;size:50;not null"`
	Title      string `gorm:"size:100;not null"`
	OrderIndex int    `gorm:"not null;index"`
}