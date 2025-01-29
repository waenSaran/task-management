package main

import (
	"log"
	"os"
	"task-management-api/config"
	"task-management-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDB()
	config.InitSupabase()

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.AuthRoutes(v1) // Authentication routes
	routes.TaskRoutes(v1) // Protected task routes

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}
