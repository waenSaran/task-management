package handlers

import (
	"encoding/json"
	"task-management-api/config"
	"task-management-api/models"
)

func AddHistory(taskID string, task models.Task) error {
	var oldTask models.Task
	if err := config.DB.First(&oldTask, "id = ?", taskID).Error; err != nil {
		return err
	}

	changes := make(map[string]map[string]string)

	if task.Title != "" && oldTask.Title != task.Title {
		changes["title"] = map[string]string{"from": oldTask.Title, "to": task.Title}
	}
	if task.Description != "" && oldTask.Description != task.Description {
		changes["description"] = map[string]string{"from": oldTask.Description, "to": task.Description}
	}
	if task.Status != "" && oldTask.Status != task.Status {
		changes["status"] = map[string]string{"from": string(oldTask.Status), "to": string(task.Status)}
	}
	if task.Assignee != nil && oldTask.Assignee != task.Assignee {
		var oldAssignee string
		var newAssignee string
		if oldTask.Assignee != nil {
			oldAssignee = *oldTask.Assignee
		}
		if task.Assignee != nil {
			newAssignee = *task.Assignee
		}
		changes["assignee"] = map[string]string{"from": oldAssignee, "to": newAssignee}
	}

	// If no changes, return nil
	if len(changes) == 0 {
		return nil
	}

	// Convert changes map to JSON
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return err
	}

	// Create a TaskHistory record
	history := models.History{
		TaskID:    taskID,
		ChangedBy: task.UpdatedBy,
		Changes:   string(changesJSON),
	}

	// Save the history
	if err := config.DB.Create(&history).Error; err != nil {
		return err
	}

	return nil
}
