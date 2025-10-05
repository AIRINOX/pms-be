package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Username string `gorm:"uniqueIndex;size:100"`
	Password string `gorm:"size:255"`
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Phone    string `gorm:"size:50"`
	RoleID   uint   `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`

	// Relationships
	Role Role `gorm:"foreignKey:RoleID"`
}

// GetKey returns the primary key of the user (required for authentication)
func (u *User) GetKey() any {
	return u.ID
}

// GetAuthIdentifierName returns the name of the unique identifier for the user
func (u *User) GetAuthIdentifierName() string {
	return "id"
}

// GetAuthIdentifier returns the unique identifier for the user
func (u *User) GetAuthIdentifier() any {
	return u.ID
}

// GetAuthPassword returns the password for the user
func (u *User) GetAuthPassword() string {
	return u.Password
}
