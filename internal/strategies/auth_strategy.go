package strategies

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/starks97/alcohol-tracker-api/internal/state"
)

type AuthStrategy interface {
	GenerateAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) ([]byte, error)
}

func NewAuthStrategy(appState *state.AppState, provider string) (AuthStrategy, error) {
	switch provider {
	case "github":
		return &GitHubStrategy{AppState: appState}, nil
	case "google":
		return &GoogleStrategy{AppState: appState}, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
