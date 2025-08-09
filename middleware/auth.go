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

		// 1️⃣ Try from HTTPOnly cookie first
		cookie := c.Cookies("token")
		if cookie != "" {
			tokenStr = cookie
		}

		// 2️⃣ Fallback to Authorization header
		if tokenStr == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 3️⃣ No token → unauthorized
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		// 4️⃣ Verify JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// 5️⃣ Redis session check
		val, err := db.RedisClient.Get(db.Ctx, tokenStr).Result()
		if err != nil || val == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired or invalid"})
		}

		// 6️⃣ Store user ID in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Locals("user_id", sub)
			}
		}

		return c.Next()
	}
}
