package routes

import (
	"task-management-api/handlers"
	"task-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CommentRoutes(route fiber.Router) {
	task := route.Group("/tasks/:id/comments", middleware.AuthMiddleware)
	comment := route.Group("/comments", middleware.AuthMiddleware)

	task.Post("/", handlers.CreateComment)
	comment.Put("/:id", handlers.UpdateComment)
	comment.Delete("/:id", handlers.DeleteComment)
}
