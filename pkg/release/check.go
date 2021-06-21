package release

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/gh"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
)

type CheckOptions struct {
	// Repository is the GitHub repository to check.
	Repository *repository.Repository

	// ReleaseName is the name of the GitHub Release to check.
	ReleaseName string

	// AssetFiles is a list of files which must all exist in the release for
	// the check to pass.
	AssetFiles []string

	// GitHubToken is the OAuth token used to talk to the GitHub API.
	GithubToken string
}

// Check checks if a GitHub repository has a Release by given name, and if the
// release contains assets with given filenames.
func Check(ctx context.Context, opts *CheckOptions) error {
	logger := hclog.FromContext(ctx).Named("release")
	gh := gh.New(ctx, opts.GithubToken)

	repo := opts.Repository

	release, resp, err := gh.Repositories.GetReleaseByTag(
		ctx, repo.Owner(), repo.Name(), opts.ReleaseName,
	)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("release %s does not exist", opts.ReleaseName)
		}

		return err
	}

	logger.Info("release exists", "release", opts.ReleaseName)

	missingMap := map[string]bool{}
	for _, filename := range opts.AssetFiles {
		filename = filepath.Base(filename)
		missingMap[filename] = true
		for _, a := range release.Assets {
			if a.GetName() == filename {
				logger.Info("asset exists", "filename", filename)
				delete(missingMap, filename)

				break
			}
		}
	}

	if len(missingMap) == 0 {
		return nil
	}

	var missing []string
	for f := range missingMap {
		missing = append(missing, f)
	}

	logger.Error("missing assets", "filenames", missing)

	return fmt.Errorf(
		"release %s is missing assets:\n- %s",
		opts.ReleaseName, strings.Join(missing, "\n-"),
	)
}
