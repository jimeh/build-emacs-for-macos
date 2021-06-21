package gh

import (
	"context"
	"os"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func New(ctx context.Context, token string) *github.Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
