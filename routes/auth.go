package routes

import (
	"task-management-api/handlers"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes sets up authentication endpoints
func AuthRoutes(route fiber.Router) {
	auth := route.Group("/auth")

	auth.Post("/signup", handlers.SignUp)
	auth.Post("/login", handlers.Login)
}
