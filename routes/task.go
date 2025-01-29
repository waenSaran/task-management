package routes

import (
	"task-management-api/handlers"
	"task-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

// TaskRoutes sets up protected task endpoints
func TaskRoutes(route fiber.Router) {
	tasks := route.Group("/tasks", middleware.AuthMiddleware) // Protect all task routes

	tasks.Post("/", handlers.CreateTask)
	tasks.Get("/:id", handlers.GetTaskById)
	tasks.Put("/:id", handlers.UpdateTask)
	tasks.Delete("/:id", handlers.DeleteTask)
}
