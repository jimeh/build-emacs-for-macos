package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v35/github"
	"github.com/urfave/cli/v2"
)

var publishCmd = &cli.Command{
	Name:      "publish",
	Usage:     "publish a release",
	UsageText: "github-release [global-options] publish [options]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "github-sha",
			Aliases: []string{"s"},
			Usage:   "Git SHA of repo to create release on",
			EnvVars: []string{"GITHUB_SHA"},
		},
		&cli.BoolFlag{
			Name:    "prerelease",
			Usage:   "Git SHA of repo to create release on",
			EnvVars: []string{"RELEASE_PRERELEASE"},
			Value:   true,
		},
	},
	Action: actionHandler(publishAction),
}

func publishAction(c *cli.Context, opts *globalOptions) error {
	gh := opts.gh
	repo := opts.repo
	plan := opts.plan

	releaseName := plan.ReleaseName()
	githubSHA := c.String("github-sha")
	prerelease := c.Bool("prerelease")

	fmt.Printf("prerelease: %+v\n", prerelease)

	assetFile, err := os.Open(plan.Archive)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching release %s\n", releaseName)

	release, resp, err := gh.Repositories.GetReleaseByTag(
		c.Context, repo.Owner, repo.Name, releaseName,
	)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("Release %s not found, creating...\n", releaseName)
			release, _, err = gh.Repositories.CreateRelease(
				c.Context, repo.Owner, repo.Name, &github.RepositoryRelease{
					Name:            &releaseName,
					TagName:         &releaseName,
					TargetCommitish: &githubSHA,
					Prerelease:      &prerelease,
				},
			)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if release.GetPrerelease() != prerelease {
		release.Prerelease = &prerelease

		release, _, err = gh.Repositories.EditRelease(
			c.Context, repo.Owner, repo.Name, release.GetID(), release,
		)
		if err != nil {
			return err
		}
	}

	assetFilename := plan.ReleaseAsset()
	fmt.Printf("Uploading asset %s\n", assetFilename)
	_, _, err = gh.Repositories.UploadReleaseAsset(
		c.Context, repo.Owner, repo.Name, release.GetID(),
		&github.UploadOptions{Name: assetFilename},
		assetFile,
	)
	if err != nil {
		return err
	}

	fmt.Printf("Release published at: %s\n", release.GetHTMLURL())

	return nil
}
