package models

import (
	"encoding/json"
	"time"
)

type History struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TaskID    string    `gorm:"type:uuid;not null" json:"task_id"`
	ChangedBy string    `gorm:"type:uuid;not null" json:"changed_by"`
	Changes   string    `gorm:"type:jsonb" json:"changes"` // Store as JSON
	ChangedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"changed_at"`

	// Relationships
	Task        Task `gorm:"foreignKey:TaskID"`
	ChangedUser User `gorm:"foreignKey:ChangedBy"`
}

type ChangeDetail struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Changes map[string]ChangeDetail

type HistoryResponse struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	ChangedBy string    `json:"changed_by"`
	Changes   Changes   `json:"changes"`
	ChangedAt time.Time `json:"changed_at"`
}

func FormatHistoryResponse(history History) (HistoryResponse, error) {
	var changes Changes

	// Unmarshal JSON string into structured Changes map
	err := json.Unmarshal([]byte(history.Changes), &changes)
	if err != nil {
		return HistoryResponse{}, err
	}

	return HistoryResponse{
		ID:        history.ID,
		TaskID:    history.TaskID,
		ChangedBy: history.ChangedBy,
		Changes:   changes,
		ChangedAt: history.ChangedAt,
	}, nil
}
