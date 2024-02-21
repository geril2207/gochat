package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get(fiber.HeaderAuthorization)
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	key := os.Getenv("JWT_KEY")
	if key == "" {
		panic("JWT_KEY not found")
	}

	// opts := jwt.ParseWithC

	return ctx.Next()
}
