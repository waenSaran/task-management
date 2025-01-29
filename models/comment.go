package models

import (
	"time"
)

type Comment struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content   string    `gorm:"not null"`
	TaskID    string    `gorm:"type:uuid;not null" json:"task_id"`
	CreatedBy string    `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:CreatedBy"`
	Task Task `gorm:"foreignKey:TaskID"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	TaskID    string    `json:"task_id"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatCommentResponse(comment Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		TaskID:    comment.TaskID,
		CreatedBy: comment.CreatedBy,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
