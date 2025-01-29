package models

import "gorm.io/gorm"

type User struct {
	// Adds created_at & updated_at automatically
	gorm.Model

	ID       string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null;default:'user'" json:"role"`
}
