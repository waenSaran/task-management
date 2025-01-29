package handlers

import (
	"task-management-api/config"
	"task-management-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserByID(c *fiber.Ctx) models.User {
	claims := c.Locals("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	var user models.User
	res := config.DB.Where("id = ?", userID).First(&user)
	if res.Error != nil {
		c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return user
}
