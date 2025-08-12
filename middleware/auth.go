package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"supa.fiber/db"
)

func Protect() fiber.Handler {
	secret := os.Getenv("SUPABASE_JWT_SECRET")

	return func(c *fiber.Ctx) error {
		var tokenStr string

		cookie := c.Cookies("token")
		if cookie != "" {
			tokenStr = cookie
		}

		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		val, err := db.RedisClient.Get(db.Ctx, tokenStr).Result()
		if err != nil || val == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired or invalid"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Locals("user_id", sub)
			}
		}

		return c.Next()
	}
}
