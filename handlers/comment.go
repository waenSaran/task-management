package handlers

import (
	"task-management-api/config"
	"task-management-api/models"
	"task-management-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateComment(c *fiber.Ctx) error {
	user := GetUserByID(c)
	taskId := c.Params("id")

	var comment models.Comment

	// validate task id
	if err := config.DB.Table("tasks").Where("id = ?", taskId).Select("id").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve task"})
	}

	// validate payload
	if err := c.BodyParser(&comment); err != nil || comment.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	comment.CreatedBy = user.ID
	comment.TaskID = taskId

	// Create the comment
	if err := config.DB.Create(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	response := models.FormatCommentResponse(comment)

	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateComment(c *fiber.Ctx) error {
	user := GetUserByID(c)
	commentID := c.Params("id")

	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve comment"})
	}

	// Check if the comment belongs to the authenticated user
	if !utils.HasPermission(comment.CreatedBy, user.ID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to update this comment"})
	}

	// Parse the update data
	var updatedComment models.Comment
	if err := c.BodyParser(&updatedComment); err != nil || updatedComment.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Update the comment
	if err := config.DB.Model(&comment).Updates(updatedComment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	response := models.FormatCommentResponse(comment)

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteComment(c *fiber.Ctx) error {
	commentID := c.Params("id")

	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve comment"})
	}

	// Check if the comment belongs to the authenticated user
	user := GetUserByID(c)
	if !utils.HasPermission(comment.CreatedBy, user.ID) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to delete this comment"})
	}

	// Delete the comment
	if err := config.DB.Delete(&comment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete comment"})
	}

	return c.Status(fiber.StatusOK).SendString("Comment deleted")
}
