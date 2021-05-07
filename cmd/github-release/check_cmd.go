package main

import (
	"fmt"
	"regexp"

	"github.com/urfave/cli/v2"
)

var nonAlphaNum = regexp.MustCompile(`[^\w-_]+`)

var checkCmd = &cli.Command{
	Name:      "check",
	Usage:     "Check if GitHub release and asset exists",
	UsageText: "github-release [global options] check [options]",
	Action:    actionHandler(checkAction),
}

func checkAction(c *cli.Context, opts *globalOptions) error {
	gh := opts.gh
	repo := opts.repo
	plan, err := LoadPlan(opts.plan)
	if err != nil {
		return err
	}

	releaseName := plan.ReleaseName()

	fmt.Printf(
		"Checking github.com/%s for release: %s\n",
		repo.String(), releaseName,
	)

	release, _, err := gh.Repositories.GetReleaseByTag(
		c.Context, repo.Owner, repo.Name, releaseName,
	)
	if err != nil {
		return err
	}

	filename := plan.ReleaseAsset()

	fmt.Printf("checking release for asset: %s\n", filename)
	for _, a := range release.Assets {
		if a.Name != nil && filename == *a.Name {
			return nil
		}
	}

	return fmt.Errorf("release does contain asset: %s", filename)
}
