package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/handlers"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

// all routes
func SetupRoutes(app *fiber.App, appState *state.AppState) {

	//app.Use(middleware.JWTAuthMiddleware())
	app.Post("/upload", handlers.UploadImageHandler)

	auth := app.Group("/auth")
	auth.Get("/google_login", handlers.GoogleLoginHandler)
	auth.Get("/google_callback", handlers.GoogleCallBack)
}
