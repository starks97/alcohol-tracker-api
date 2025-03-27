package strategies

import (
	"context"
	"io"

	"golang.org/x/oauth2"

	"github.com/starks97/alcohol-tracker-api/internal/state"
)

type GoogleStrategy struct {
	AppState *state.AppState
}

func (g *GoogleStrategy) GenerateAuthURL(state string) string {
	return g.AppState.Config.GoogleLoginConfig.AuthCodeURL(state)
}

func (g *GoogleStrategy) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.AppState.Config.GoogleLoginConfig.Exchange(ctx, code)
}

func (g *GoogleStrategy) GetUserInfo(token *oauth2.Token) ([]byte, error) {
	// Directly using HttpClient.Get instead of creating a new request
	resp, err := g.AppState.HttpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Set the Authorization header for the request
	resp.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Read and return the response body
	return io.ReadAll(resp.Body)

}
