package release

import (
	"context"
	"regexp"

	"github.com/google/go-github/v35/github"
	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
)

type BulkOptions struct {
	Repository  *repository.Repository
	NamePattern string
	Prerelease  *bool
	DryRun      bool
	GithubToken string
}

func Bulk(ctx context.Context, opts *BulkOptions) error {
	logger := hclog.FromContext(ctx).Named("release")
	gh := gh.New(ctx, opts.GithubToken)

	nameMatcher, err := regexp.Compile(opts.NamePattern)
	if err != nil {
		return err
	}

	nextPage := 1
	lastPage := 1

	for nextPage <= lastPage {
		releases, resp, err := gh.Repositories.ListReleases(
			ctx, opts.Repository.Owner(), opts.Repository.Name(),
			&github.ListOptions{
				Page:    nextPage,
				PerPage: 100,
			},
		)
		if err != nil {
			return err
		}

		nextPage = resp.NextPage
		lastPage = resp.LastPage

		for _, r := range releases {
			if !nameMatcher.MatchString(r.GetName()) {
				continue
			}

			logger.Info("match found", "release", r.GetName())

			var changes []interface{}
			if opts.Prerelease != nil && r.GetPrerelease() != *opts.Prerelease {
				changes = append(changes, "prerelease", *opts.Prerelease)
				r.Prerelease = opts.Prerelease
			}

			if len(changes) > 0 {
				changes = append(
					[]interface{}{"release", r.GetName()}, changes...,
				)
				logger.Info("modifying", changes...)
				if !opts.DryRun {
					r, _, err = gh.Repositories.EditRelease(
						ctx, opts.Repository.Owner(), opts.Repository.Name(),
						r.GetID(), r,
					)
					if err != nil {
						return err
					}
				}
			}
		}

		if nextPage == 0 || lastPage == 0 {
			break
		}
	}

	return nil
}
