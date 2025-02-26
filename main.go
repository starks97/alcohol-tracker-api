package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/handlers"
)

func main() {
	app := fiber.New()

	// Ruta de prueba
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong"})
	})

	// Configurar rutas
	handlers.SetupRoutes(app)

	// Iniciar servidor
	port := "8080"
	fmt.Println("🚀 Servidor corriendo en http://localhost:" + port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Error iniciando el servidor:", err)
	}
}
