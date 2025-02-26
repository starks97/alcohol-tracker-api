package handlers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/internal/services"
)

func UploadImageHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("Error to find the error: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "File not found.",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		log.Printf("Error to open the file: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "It was not possible to process the file",
		})

	}

	defer fileContent.Close()

	processedImage, err := services.ProcessImage(fileContent)
	if err != nil {
		log.Printf("Error to process the image: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing image",
		})
	}

	println(processedImage)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "The image was processed and sent successfully to the NN",
		"image":   processedImage,
	})
}
