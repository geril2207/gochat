package auth

import "github.com/gofiber/fiber/v2"

func AuthRoutes(router fiber.Router) fiber.Router {
	authGroup := router.Group("/auth")

	authGroup.Post("/login", Login)
	authGroup.Post("/register", Register)

	return authGroup
}
