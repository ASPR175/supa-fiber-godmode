package middleware

import (
	"time"

	"supa.fiber/db"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func BlacklistToken(token string, duration time.Duration) error {
	return db.RedisClient.Set(db.Ctx, token, "blacklisted", duration).Err()
}

func IsTokenBlacklisted(token string) (bool, error) {
	result, err := db.RedisClient.Get(db.Ctx, token).Result()

	if err == redis.Nil {

		return false, nil
	}
	if err != nil {
		return false, err
	}

	return result == "blacklisted", nil
}
func ExtractTokenFromHandler(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if len(authHeader) >= 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
func TokenBlacklistMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := ExtractTokenFromHandler(c)

		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Message": "Empty Token",
			})
		}
		isBlacklisted, err := IsTokenBlacklisted(token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Message": "Error occured to Get Token",
			})
		}
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has been revoked",
			})
		}
		return c.Next()
	}
}
