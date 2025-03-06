package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/routes"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	routes.SetupRoutes(app)

	port := "8080"
	fmt.Println("🚀 Server running on http://localhost:" + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
