package authen

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/starks97/alcohol-tracker-api/internal/dtos"
	"github.com/starks97/alcohol-tracker-api/internal/entities"
	"github.com/starks97/alcohol-tracker-api/internal/exceptions"
	"github.com/starks97/alcohol-tracker-api/internal/repositories"
	"github.com/starks97/alcohol-tracker-api/internal/responses"
	"github.com/starks97/alcohol-tracker-api/internal/state"
	"github.com/starks97/alcohol-tracker-api/internal/strategies"
	"github.com/starks97/alcohol-tracker-api/internal/utils"
)

// OAuthLoginHandler initiates the Google OAuth2 login flow.
// It retrieves the Google OAuth2 configuration from the Fiber context,
// generates an authorization URL, and redirects the user to Google's login page.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context.
//
// Returns:
//   - error: An error if the redirection fails, or nil if successful.
func OAuthLoginHandler(c *fiber.Ctx) error {
	appState := c.Locals("appState").(*state.AppState)
	provider := c.Params("provider")

	authStrategy, err := strategies.NewAuthStrategy(appState, provider)
	if err != nil {
		log.Printf("Error creating auth strategy: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}
	if authStrategy == nil {
		log.Println("AuthStrategy is nil")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid auth strategy"})
	}

	state, err := utils.GenerateRandomString(32)
	if err != nil {
		log.Print("can not generate random state", err)
		return fmt.Errorf("GenerateRandomString: %w", exceptions.HandlerErrorResponse(c, exceptions.ErrTokenNotGenerated))
	}

	cookie := fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "localhost",
	}

	c.Cookie(&cookie)

	url := authStrategy.GenerateAuthURL(state)

	// Generate the authorization URL using a fixed state value.

	// Redirect the user to the generated URL with a "See Other" status.
	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)

	// Return the URL as JSON (though the redirect should take precedence).
	return c.JSON(url)
}

// OAuthCallBackHandler handles the callback from Google's OAuth 2.0 authorization server.
// It verifies the state parameter to prevent CSRF attacks, exchanges the authorization code
// for an access token, retrieves user information from Google, and either creates a new user
// or logs in an existing user. Finally, it generates JWT tokens and sets a refresh token cookie.
//
// Parameters:
//   - c: *fiber.Ctx - The Fiber context.
//
// Returns:
//   - error: An error if any step of the process fails, or nil if successful.
func OAuthCallBackHandler(c *fiber.Ctx) error {
	// Retrieve application state and context from Fiber locals.
	appState := c.Locals("appState").(*state.AppState)
	ctx := c.Locals("ctx").(context.Context)
	provider := c.Params("provider")

	tokenService := utils.NewTokenService(appState)

	authStrategy, err := strategies.NewAuthStrategy(appState, provider)
	if err != nil {
		return exceptions.HandlerErrorResponse(c, err)
	}

	// Initialize user repository.
	userRepo := repositories.NewUserRepository(appState.DB)

	// Retrieve the state parameter from the cookie to prevent CSRF attacks.
	cookieState := c.Cookies("oauth_state")

	// Retrieve the authorization code and state from the query parameters.
	code := c.Query("code")
	queryState := c.Query("state")

	// Verify that the state from the cookie matches the state from the query parameters.
	if cookieState != queryState {
		return c.SendString("States don't Match!!")
	}

	// Exchange the authorization code for an access token from Google.
	token, err := authStrategy.ExchangeCode(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange token:", err)
		return exceptions.HandlerErrorResponse(c, exceptions.ErrExchangeToken)
	}

	// Retrieve user information from Google's userinfo endpoint using the access token.
	userData, err := authStrategy.GetUserInfo(token)
	if err != nil {
		log.Println("Failed to get user info:", err)
		return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotFound)
	}

	// Unmarshal the user information into a GoogleUser struct.
	var oauthUser dtos.OAuthDto
	err = json.Unmarshal(userData, &oauthUser)
	if err != nil {
		log.Println("Failed to unmarshal user info:", err)
		return exceptions.HandlerErrorResponse(c, exceptions.ErrToUnmarshalUserInfo)
	}

	// Check if the user exists in the database.
	user, err := userRepo.GetUserByEmail(oauthUser.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &entities.User{
				Name:           oauthUser.Name,
				Email:          oauthUser.Email,
				Provider:       &provider,
				ProviderID:     &oauthUser.ID,
				ProfilePicture: &oauthUser.Picture,
			}
			_, err = userRepo.CreateUser(user)
			if err != nil {
				log.Println("Failed to create user:", err)
				return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotCreated)
			}
		} else {
			log.Println("Failed to get user:", err)
			return exceptions.HandlerErrorResponse(c, fmt.Errorf("failed to get user: %w", err))
		}
	} else {
		//user exist update user
		user.Provider = &provider
		user.ProviderID = &oauthUser.ID
		user.ProfilePicture = &oauthUser.Picture
		user.Name = oauthUser.Name
		user.ProviderRefreshToken = &token.AccessToken

		_, err = userRepo.UpdateUser(user)
		if err != nil {
			log.Println("Failed to update user:", err)
			return exceptions.HandlerErrorResponse(c, exceptions.ErrUserNotUpdated)
		}
	}

	// Generate and store JWT tokens in Redis, and set a refresh token cookie.
	tokenResult, err := tokenService.StoreToken(c, ctx, user.ID, "both")
	if err != nil {
		return err
	}

	// Create a login response with the access token.
	accessToken := responses.LoginResponse{
		AccessToken: *tokenResult.Token,
	}

	// Return a success response with the access token.
	return c.JSON(responses.SuccessResponse{
		Status: "success",
		Data:   accessToken,
	})
}
