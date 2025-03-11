package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/database"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/routes"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return exceptions.NewCustomErrorResponse(ctx, err)
		},
	})
	//helps with context
	ctx := context.Background()
	httpClient := &http.Client{}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	redisClient, err := database.NewRedisClient(cfg, ctx)
	if err != nil {
		log.Fatalf("Error initializing Redis client: %v", err)
	}

	db := database.ConnectDB(cfg)

	appState := &state.AppState{
		DB:         db,
		Redis:      redisClient,
		Config:     cfg,
		HttpClient: httpClient,
	}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("appState", appState)
		c.Locals("ctx", ctx)
		return c.Next()
	})

	app.Use(cors.New(
		cors.Config{
			AllowOrigins:     appState.Config.ClientOrigin,
			AllowHeaders:     "Authorization, Content-Type, Accept, Access-Control-Allow-Origin",
			MaxAge:           3600,
			AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
			AllowCredentials: true,
		},
	))

	routes.SetupRoutes(app, appState)

	port := "8080"
	fmt.Println("🚀 Server running on http://localhost:" + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
