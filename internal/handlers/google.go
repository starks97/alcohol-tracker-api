package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/models"
)

func GoogleLoginHandler(c *fiber.Ctx) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	url := cfg.GoogleLoginConfig.AuthCodeURL("randomstate")
	c.Status((fiber.StatusSeeOther))
	c.Redirect(url)
	return c.JSON(url)
}

func GoogleCallBack(c *fiber.Ctx) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	state := c.Query("state")
	if state != "randomstate" {
		return c.SendString("States don't Match!!")
	}

	code := c.Query("code")

	token, err := cfg.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to exchange token",
		})
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println("Failed to get user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user info",
		})
	}
	defer res.Body.Close()

	userData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read user info",
		})
	}

	var googleUser models.GoogleUser
	err = json.Unmarshal(userData, &googleUser)
	if err != nil {
		log.Println("Failed to unmarshal user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unmarshal user info",
		})
	}

	return c.JSON(fiber.Map{
		"name":        googleUser.Name,
		"email":       googleUser.Email,
		"picture":     googleUser.VerifiedEmail,
		"provider_id": googleUser.ID,
		"provider":    "google",
		"token":       token.AccessToken,
	})
}
