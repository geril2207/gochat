package router

import (
	"github.com/geril2207/gochat/apps/server/auth"
	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app *fiber.App) fiber.Router {
	api := app.Group("/api")

	auth.AuthRoutes(api)

	return api
}
