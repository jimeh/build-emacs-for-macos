package main

import (
	"context"
	"net/http"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(ctx context.Context, token string) *github.Client {
	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	return github.NewClient(tc)
}
