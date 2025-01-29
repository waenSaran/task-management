package handlers

import (
	"task-management-api/config"
	"task-management-api/models"
	"task-management-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllTasks fetches all tasks with optional filters
func GetAllTasks(c *fiber.Ctx) error {
	var tasks []models.Task

	// Initialize query builder
	query := config.DB.Model(&models.Task{})

	// Apply filters based on query params
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	createdBy := c.Query("createdBy")
	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	title := c.Query("title")
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// Execute query and fetch tasks
	if err := query.Find(&tasks).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No tasks found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
	}

	return c.JSON(tasks)
}

func CreateTask(c *fiber.Ctx) error {
	user := GetUserByID(c)

	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	task.CreatedBy = user.ID
	task.UpdatedBy = user.ID

	// Create the task
	if err := config.DB.Preload("CreatedUser").Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	response := models.FormatTaskResponse(task)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func GetTaskById(c *fiber.Ctx) error {
	taskID := c.Params("id")

	var task models.Task
	if err := config.DB.Preload("CreatedUser").Preload("UpdatedUser").Preload("AssigneeUser").First(&task, "id = ?", taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve task"})
	}

	res := models.FormatTaskResponse(task)

	return c.JSON(res)
}

func UpdateTask(c *fiber.Ctx) error {
	user := GetUserByID(c)
	taskID := c.Params("id")

	var task models.Task
	if err := config.DB.First(&task, "id = ?", taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve task"})
	}

	// Parse the update data
	var updatedTask models.Task
	if err := c.BodyParser(&updatedTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// validate status if status exists
	if updatedTask.Status != "" && !validateStatus(string(updatedTask.Status)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Status must be one of: TODO, IN_PROGRESS, DONE"})
	}

	payload := models.Task{
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      updatedTask.Status,
		Assignee:    updatedTask.Assignee,
		UpdatedBy:   user.ID,
	}

	// Update history before updating the task
	if err := AddHistory(taskID, payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update history"})
	}

	// Update the task
	if err := config.DB.Model(&task).Updates(payload).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	response := models.FormatTaskResponse(task)

	return c.JSON(response)
}

func DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")

	var task models.Task
	if err := config.DB.First(&task, "id = ?", taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve task"})
	}

	// Check if the task belongs to the authenticated user
	user := GetUserByID(c)
	if !utils.HasPermission(task.CreatedBy, user.ID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to delete this task"})
	}

	// Delete associated comments
	if err := config.DB.Table("comments").Where("task_id = ?", taskID).Delete(&models.Comment{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comments before deleting task"})
	}

	// Delete associated histories
	if err := config.DB.Table("histories").Where("task_id = ?", taskID).Delete(&models.History{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete history before deleting task"})
	}

	// Delete the task
	if err := config.DB.Delete(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.Status(fiber.StatusOK).SendString("Task deleted")
}

// return true if stutus is in type enum utils.Status
func validateStatus(status string) bool {
	return status == string(utils.Todo) || status == string(utils.InProgress) || status == string(utils.Done)
}
