package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"supa.fiber/db"
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

	// Middlewares
	app.Use(helmet.New())
	app.Use(logger.New())

	// CORS so frontend (3000) can talk to backend (8080)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Routes
	routes.UserRoutes(app, queries)
	routes.TaskRoutes(app, queries)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // backend port
	}
	log.Fatal(app.Listen(":" + port))
}
