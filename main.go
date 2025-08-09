package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"supa.fiber/db"

	// "supa.fiber/handlers"

	"supa.fiber/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load godotenv")
	}

	conn := db.ConnectDB()
	defer conn.Close()

	queries := db.New(conn)
	db.InitRedis()

	app := fiber.New()

	app.Use(helmet.New())
	app.Use(logger.New())

	// Public routes
	routes.UserRoutes(app, queries)

	// Protected group

	routes.TaskRoutes(app, queries)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
