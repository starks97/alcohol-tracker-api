package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/database"
	"github.com/starks97/alcohol-tracker-api/internal/routes"
	"github.com/starks97/alcohol-tracker-api/internal/state"
)

func main() {
	app := fiber.New()
	//helps with context
	ctx := context.Background()

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
		DB:     db,
		Redis:  redisClient,
		Config: cfg,
	}
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("appState", appState)
		c.Locals("ctx", ctx)
		return c.Next()
	})

	routes.SetupRoutes(app, appState)

	port := "8080"
	fmt.Println("🚀 Server running on http://localhost:" + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
