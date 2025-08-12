package handlers

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	// "github.com/redis/go-redis/v9"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
	"supa.fiber/db"
	"supa.fiber/utils"
)

type UserHandler struct {
	Q db.Querier
}

func NewUserHandler(q db.Querier) *UserHandler {
	return &UserHandler{Q: q}
}

func (h *UserHandler) SignUp(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user, err := h.Q.CreateUser(context.Background(), db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashed),
	})
	if err != nil {
		log.Printf("CreateUser:%+v", err)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token creation failed"})
	}

	err = db.RedisClient.Set(db.Ctx, token, user.ID.String(), 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		"token":   token,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.Q.GetUserByEmail(context.Background(), req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token creation failed"})
	}

	err = db.RedisClient.Set(db.Ctx, token, user.ID.String(), 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func Logout(c *fiber.Ctx) error {

	token := c.Cookies("token")

	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > len("Bearer ") {
			token = authHeader[len("Bearer "):]
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	exists, err := db.RedisClient.Exists(db.Ctx, token).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking session"})
	}

	if exists == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token invalid or already logged out"})
	}

	if err := db.RedisClient.Del(db.Ctx, token).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Logout failed"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"message": "Logged out"})
}
