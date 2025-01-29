package handlers

import (
	"task-management-api/config"
	"task-management-api/models"
	"task-management-api/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "email and password are required",
		})
	}

	user := models.User{
		Email:    req.Email,
		Password: utils.GeneratePassword(req.Password),
	}
	res := config.DB.Create(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": res.Error.Error(),
		})
	}

	// Login user automatically
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Successfully signed up",
		"token":   token,
	})
}

// Login route
func Login(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "email and password are required",
		})
	}

	var user models.User
	res := config.DB.Where("email = ?", req.Email).First(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	if !utils.ComparePassword(user.Password, req.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
