package models

import "gorm.io/gorm"

type Comment struct {
	// Adds created_at & updated_at automatically
	gorm.Model

	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content   string `gorm:"not null"`
	TaskID    string `gorm:"type:uuid;not null"`
	CreatedBy string `gorm:"type:uuid" json:"created_by"`

	// Relationships
	User User `gorm:"foreignKey:CreatedBy"`
	Task Task `gorm:"foreignKey:TaskID"`
}
