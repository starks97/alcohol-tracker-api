package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/handlers/authen"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

// all routes
func SetupRoutes(app *fiber.App, appState *state.AppState) {

	//app.Use(middleware.JWTAuthMiddleware())

	auth := app.Group("/auth")
	auth.Get("/:provider", authen.OAuthLoginHandler)
	auth.Get("/:provider/callback", authen.OAuthCallBackHandler)
	auth.Post("/register", authen.Register)
	auth.Post("/login", authen.LoginHandler)
}
