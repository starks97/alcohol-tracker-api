package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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
			return exceptions.HandlerErrorResponse(ctx, err)
		},
	})

	//helps with context
	ctx := context.Background()

	//http
	httpClient := &http.Client{}

	//load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	//redis client
	redisClient, err := database.NewRedisClient(cfg, ctx)
	if err != nil {
		log.Fatalf("Error initializing Redis client: %v", err)
	}

	//database connection
	db := database.ConnectDB(cfg)

	//validator
	validator := exceptions.Init()

	//initialize state
	appState := &state.AppState{
		DB:         db,
		Redis:      redisClient,
		Config:     cfg,
		HttpClient: httpClient,
		Validator:  validator,
	}

	//set interfaces available to routes
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("appState", appState)
		c.Locals("ctx", ctx)
		return c.Next()
	}, cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(cfg.ClientOrigin, origin)
		},
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization",
		AllowCredentials: true,
	}))

	//pass params to routes
	routes.SetupRoutes(app, appState)

	port := "8080"
	fmt.Println("🚀 Server running on http://localhost:" + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
