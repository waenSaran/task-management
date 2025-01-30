package models

import (
	"task-management-api/utils"
	"time"
)

type TaskStatus string

const (
	Todo       TaskStatus = "TODO"
	InProgress TaskStatus = "IN_PROGRESS"
	InReview   TaskStatus = "IN_REVIEW"
	Done       TaskStatus = "DONE"
)

type Task struct {
	ID          string       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      utils.Status `gorm:"type:varchar(20);default:'TODO'" json:"status"`
	Assignee    *string      `gorm:"type:uuid;default:NULL" json:"assignee"`
	CreatedBy   string       `gorm:"type:uuid" json:"created_by"`
	UpdatedBy   string       `gorm:"type:uuid" json:"updated_by"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	CreatedUser  User `gorm:"foreignKey:CreatedBy" json:"created_user"`
	UpdatedUser  User `gorm:"foreignKey:UpdatedBy" json:"updated_user"`
	AssigneeUser User `gorm:"foreignKey:Assignee" json:"assignee_user"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Assignee    *string   `json:"assignee,omitempty"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskDetailsResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Assignee    *string           `json:"assignee,omitempty"`
	CreatedBy   string            `json:"created_by"`
	UpdatedBy   string            `json:"updated_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Comments    []CommentResponse `json:"comments"`
	History     []HistoryResponse `json:"history"`
}

// Function to convert Task model to response format
func FormatTaskResponse(task Task) TaskResponse {
	// using user id instead of email for task creation
	createdBy := task.CreatedBy
	updatedBy := task.UpdatedBy
	assignee := task.Assignee
	if task.CreatedUser.Email != "" {
		createdBy = task.CreatedUser.Email
	}
	if task.UpdatedUser.Email != "" {
		updatedBy = task.UpdatedUser.Email
	}
	if task.AssigneeUser.Email != "" {
		*assignee = task.AssigneeUser.Email
	}

	return TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Assignee:    assignee,
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
