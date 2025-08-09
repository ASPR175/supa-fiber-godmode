package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"supa.fiber/db"
)

type TaskHandler struct {
	Q db.Querier
}

func NewTaskHandler(q db.Querier) *TaskHandler {
	return &TaskHandler{Q: q}
}
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid request"})
	}

	uidStr := c.Locals("user_id").(string)
	userID := pgtype.UUID{}
	if err := userID.Scan(uidStr); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid UserID"})
	}

	task, err := h.Q.CreateTask(context.Background(), db.CreateTaskParams{
		Title: req.Title,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		UserID: userID,
	})
	if err != nil {
		log.Println("CreateTask error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}
	return c.JSON(task)
}

func (h *TaskHandler) GetAllMyTasks(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	userID := pgtype.UUID{}
	if err := userID.Scan(uidStr); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid UserID"})
	}

	tasks, err := h.Q.GetTasksByUserId(context.Background(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"Message": "Couldn't fetch"})
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	var req struct {
		TaskID string `json:"task_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid request"})
	}

	taskID := pgtype.UUID{}
	if err := taskID.Scan(req.TaskID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid TaskID"})
	}

	uidStr := c.Locals("user_id").(string)
	userID := pgtype.UUID{}
	if err := userID.Scan(uidStr); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"Message": "Invalid UserID"})
	}

	rowsAffected, err := h.Q.DeleteTask(context.Background(), db.DeleteTaskParams{
		ID:     taskID,
		UserID: userID,
	})
	if err != nil {
		log.Println("DeleteTask error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete task"})
	}

	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"Message": "Task not found or you are not the owner"})
	}

	return c.JSON(fiber.Map{"Message": "Deleted Task"})

}
