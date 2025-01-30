package routes

import (
	"task-management-api/handlers"
	"task-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func TaskRoutes(route fiber.Router) {
	publicRoute := route.Group("/tasks")
	withAuthRoute := route.Group("/tasks", middleware.AuthMiddleware)

	publicRoute.Get("/", handlers.GetAllTasks)
	publicRoute.Get("/:id", handlers.GetTaskById)

	// Only authenticated users can create, update, and delete tasks
	withAuthRoute.Post("/", handlers.CreateTask)
	withAuthRoute.Put("/:id", handlers.UpdateTask)
	withAuthRoute.Delete("/:id", handlers.DeleteTask)
}
