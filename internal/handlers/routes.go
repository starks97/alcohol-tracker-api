package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// all routes
func SetupRoutes(app *fiber.App) {
	app.Post("/upload", UploadImageHandler)
}
