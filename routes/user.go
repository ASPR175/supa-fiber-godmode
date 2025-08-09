package routes

import (
	"github.com/gofiber/fiber/v2"
	"supa.fiber/db"
	"supa.fiber/handlers"
)

func UserRoutes(app *fiber.App, q db.Querier) {
	handler := handlers.NewUserHandler(q)
	user := app.Group("/api/user")
	user.Post("/signup", handler.SignUp)
	user.Post("/login", handler.Login)
	user.Post("/logout", handlers.Logout)
}
