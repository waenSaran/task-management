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

	pageQuery := c.Query("page", "1")          // Default to page 1
	pageSizeQuery := c.Query("pageSize", "10") // Default to pageSize 10
	// parse page and pageSize
	page := utils.StrToInt(pageQuery)
	pageSize := utils.StrToInt(pageSizeQuery)
	title := c.Query("title", "")
	status := c.Query("status", "")
	assignee := c.Query("assignee", "")
	createdBy := c.Query("createdBy", "")

	// Initialize query builder
	query := config.DB.Model(&models.Task{})

	// Apply filters based on query params
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if assignee != "" {
		query = query.Where("assignee = ?", assignee)
	}

	var count int64
	query.Count(&count)

	var totalPages int
	if count%int64(pageSize) == 0 {
		totalPages = int(count) / pageSize
	} else {
		totalPages = int(count)/pageSize + 1
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Execute query and fetch tasks
	if err := query.Find(&tasks).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No tasks found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve tasks"})
	}

	// format each task and get total count
	var data []interface{}
	for _, task := range tasks {
		data = append(data, models.FormatTaskResponse(task))
	}

	response := models.TransformPagination(&models.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      int(count),
		TotalPages: totalPages,
		Data:       data,
	})

	return c.JSON(response)
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

	taskDetails, err := getTaskWithDetails(taskID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve task"})
	}

	return c.JSON(taskDetails)
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
		UpdatedBy:   user.ID,
	}

	updatedTask.UpdatedBy = user.ID // for history

	// Update history before updating the task
	if err := AddHistory(taskID, updatedTask); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update history"})
	}

	// validate assignee if assignee exists in payload
	if updatedTask.Assignee != nil {
		// remove assignee if assignee from payload is empty string
		if *updatedTask.Assignee == "" {
			if err := config.DB.Model(&task).Update("assignee", nil).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update assignee to null"})
			}
		} else if !validateAssignee(*updatedTask.Assignee) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Assignee must be a valid user ID"})
		} else {
			payload.Assignee = updatedTask.Assignee
		}
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
	return status == string(utils.Todo) || status == string(utils.InProgress) || status == string(utils.Done) || status == string(utils.Archive)
}

func getTaskWithDetails(taskID string) (models.TaskDetailsResponse, error) {
	var task models.Task
	var comments []models.Comment
	var history []models.History

	// Fetch the task by ID
	if err := config.DB.Preload("CreatedUser").Preload("UpdatedUser").Preload("AssigneeUser").First(&task, "id = ?", taskID).Error; err != nil {
		return models.TaskDetailsResponse{}, err
	}

	// Fetch the comments for the task
	if err := config.DB.Preload("User").Where("task_id = ?", taskID).Find(&comments).Error; err != nil {
		return models.TaskDetailsResponse{}, err
	}

	// Fetch the task history
	if err := config.DB.Preload("ChangedUser").Where("task_id = ?", taskID).Find(&history).Error; err != nil {
		return models.TaskDetailsResponse{}, err
	}

	var commentsResponse []models.CommentResponse
	var historyResponse []models.HistoryResponse
	// Format the comments
	for _, comment := range comments {
		commentsResponse = append(commentsResponse, models.FormatCommentResponse(comment))
	}
	// Format the history
	for _, history := range history {
		h, err := models.FormatHistoryResponse(history)
		if err != nil {
			return models.TaskDetailsResponse{}, err
		}
		historyResponse = append(historyResponse, h)
	}

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
	// Format the response
	response := models.TaskDetailsResponse{
		ID:          task.ID,
		Title:       task.Title,
		Status:      string(task.Status),
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
		Assignee:    assignee,
		Comments:    commentsResponse,
		History:     historyResponse,
	}

	return response, nil
}

func validateAssignee(userID string) bool {
	if err := config.DB.First(&models.User{}, "id = ?", userID).Error; err != nil {
		return false
	}
	return true
}
