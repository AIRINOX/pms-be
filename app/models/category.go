package models

import (
	"github.com/goravel/framework/database/orm"
)

type Category struct {
	orm.Model
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	ParentID    *uint  `gorm:"index"`

	// Relationships
	Parent   *Category  `gorm:"foreignKey:ParentID"`
	Children []Category `gorm:"foreignKey:ParentID"`
	Articles []Article  `gorm:"foreignKey:CategoryID"`
}
