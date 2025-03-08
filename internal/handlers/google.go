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

// GoogleLoginHandler initiates the Google OAuth2 login flow.
// It retrieves the Google OAuth2 configuration from the Fiber context,
// generates an authorization URL, and redirects the user to Google's login page.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context.
//
// Returns:
//   - error: An error if the redirection fails, or nil if successful.
func GoogleLoginHandler(c *fiber.Ctx) error {
	cfg := c.Locals("cfg").(*config.Config)

	// Generate the authorization URL using a fixed state value.
	url := cfg.GoogleLoginConfig.AuthCodeURL("randomstate")

	// Redirect the user to the generated URL with a "See Other" status.
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)

	// Return the URL as JSON (though the redirect should take precedence).
	return c.JSON(url)
}

// GoogleCallBack handles the callback from Google OAuth2 after the user authorizes the application.
// It verifies the state parameter, exchanges the authorization code for an access token,
// retrieves user information from Google's userinfo endpoint, and returns the user data as JSON.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context.
//
// Returns:
//   - error: An error if any step fails (state verification, token exchange, user info retrieval, unmarshaling),
//     or nil if successful.
func GoogleCallBack(c *fiber.Ctx) error {
	cfg := c.Locals("cfg").(*config.Config)

	// Verify the state parameter to prevent CSRF attacks.
	state := c.Query("state")
	if state != "randomstate" {
		return c.SendString("States don't Match!!")
	}

	// Extract the authorization code from the query parameters.
	code := c.Query("code")

	// Exchange the authorization code for an access token.
	token, err := cfg.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "Failed to exchange token",
		})
	}

	// Retrieve user information from Google's userinfo endpoint using the access token.
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Println("Failed to get user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "Failed to get user info",
		})
	}
	defer res.Body.Close()

	// Read the user information from the response body.
	userData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "Failed to read user info",
		})
	}

	// Unmarshal the user information into a GoogleUser struct.
	var googleUser models.GoogleUser
	err = json.Unmarshal(userData, &googleUser)
	if err != nil {
		log.Println("Failed to unmarshal user info:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failed",
			"message": "Failed to parse user info",
		})
	}

	// Return the user data as JSON.
	return c.JSON(fiber.Map{
		"name":        googleUser.Name,
		"email":       googleUser.Email,
		"picture":     googleUser.VerifiedEmail,
		"provider_id": googleUser.ID,
		"provider":    "google",
		"token":       token.AccessToken,
	})
}
