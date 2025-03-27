package strategies

import (
	"context"
	"io"

	"golang.org/x/oauth2"

	"github.com/starks97/alcohol-tracker-api/internal/state"
)

type GitHubStrategy struct {
	AppState *state.AppState
}

func (git *GitHubStrategy) GenerateAuthURL(state string) string {
	return git.AppState.Config.GithubLoginConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (git *GitHubStrategy) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return git.AppState.Config.GithubLoginConfig.Exchange(ctx, code)
}

func (git *GitHubStrategy) GetUserInfo(token *oauth2.Token) ([]byte, error) {
	resp, err := git.AppState.HttpClient.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
