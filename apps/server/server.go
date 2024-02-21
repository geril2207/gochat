package main

import (
	"errors"
	"log"

	_ "github.com/geril2207/gochat/apps/server/docs"
	"github.com/geril2207/gochat/apps/server/router"
	"github.com/geril2207/gochat/packages/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title GoChat Api Documentation
// @description This is a simpe chat server implementation in golang
// @version 1.0
// @BasePath /api
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = db.DbConnect()
	if err != nil {
		panic(err)
	}
	defer db.DbCloseConnection()
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})
	app.Use(logger.New(logger.Config{
		Format: "${latency} ${path} ${method} ${status} ${error}\n",
	}))
	app.Use(recover.New())

	app.Get("/docs/*", swagger.HandlerDefault)
	router.SetupApiRoutes(app)

	app.Listen(":4000")
}

func customErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		code = fiberError.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return status code with error message
	return ctx.Status(code).JSON(fiber.Map{
		"message": err.Error(),
	})
}
