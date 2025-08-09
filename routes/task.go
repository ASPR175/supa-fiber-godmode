package routes

import (
	"github.com/gofiber/fiber/v2"
	"supa.fiber/db"
	"supa.fiber/handlers"
	"supa.fiber/middleware"
)

func TaskRoutes(app *fiber.App, q db.Querier) {
	handler := handlers.NewTaskHandler(q)
	task := app.Group("/api/task", middleware.Protect())
	task.Post("/createtask", handler.CreateTask)
	task.Post("/getall", handler.GetAllMyTasks)
	task.Post("/deletetask", handler.DeleteTask)
}
