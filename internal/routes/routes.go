package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/handlers"
)

// all routes
func SetupRoutes(app *fiber.App) {
	app.Post("/upload", handlers.UploadImageHandler)

	auth := app.Group("/auth")
	auth.Get("/google_login", handlers.GoogleLoginHandler)
	auth.Get("/google_callback", handlers.GoogleCallBack)
}
