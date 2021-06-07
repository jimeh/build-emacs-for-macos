package cli

import (
	"os"

	"github.com/jimeh/build-emacs-for-macos/pkg/plan"
	"github.com/jimeh/build-emacs-for-macos/pkg/release"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
	cli2 "github.com/urfave/cli/v2"
)

type releaseOptions struct {
	Plan        *plan.Plan
	Repository  *repository.Repository
	Name        string
	GithubToken string
}

func releaseCmd() *cli2.Command {
	tokenDefaultText := ""
	if len(os.Getenv("GITHUB_TOKEN")) > 0 {
		tokenDefaultText = "***"
	}

	return &cli2.Command{
		Name:  "release",
		Usage: "manage GitHub releases",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name: "plan",
				Usage: "path to build plan YAML file produced by " +
					"emacs-builder plan",
				Aliases:   []string{"p"},
				EnvVars:   []string{"EMACS_BUILDER_PLAN"},
				TakesFile: true,
			},
			&cli2.StringFlag{
				Name:    "repository",
				Aliases: []string{"repo", "r"},
				Usage: "owner/name of GitHub repo to check for release, " +
					"ignored if a plan is provided",
				EnvVars: []string{"GITHUB_REPOSITORY"},
				Value:   "jimeh/emacs-builds",
			},
			&cli2.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage: "name of release to operate on, ignored if plan " +
					"is provided",
			},
			&cli2.StringFlag{
				Name:        "github-token",
				Usage:       "GitHub API Token",
				EnvVars:     []string{"GITHUB_TOKEN"},
				DefaultText: tokenDefaultText,
			},
		},
		Subcommands: []*cli2.Command{
			releaseCheckCmd(),
		},
	}
}

func releaseActionWrapper(
	f func(*cli2.Context, *Options, *releaseOptions) error,
) func(*cli2.Context) error {
	return actionWrapper(func(c *cli2.Context, opts *Options) error {
		rOpts := &releaseOptions{
			Name:        c.String("name"),
			GithubToken: c.String("github-token"),
		}

		if r := c.String("repository"); r != "" {
			var err error
			rOpts.Repository, err = repository.NewGitHub(r)
			if err != nil {
				return err
			}
		}

		if f := c.String("plan"); f != "" {
			p, err := plan.Load(f)
			if err != nil {
				return err
			}

			rOpts.Plan = p
		}

		return f(c, opts, rOpts)
	})
}

func releaseCheckCmd() *cli2.Command {
	return &cli2.Command{
		Name: "check",
		Usage: "check if a GitHub release exists and has specified " +
			"asset files",
		ArgsUsage: "[<asset-file> ...]",
		Action:    releaseActionWrapper(releaseCheckAction),
	}
}

func releaseCheckAction(
	c *cli2.Context,
	opts *Options,
	rOpts *releaseOptions,
) error {
	rlsOpts := &release.CheckOptions{
		Repository:  rOpts.Repository,
		ReleaseName: rOpts.Name,
		AssetFiles:  c.Args().Slice(),
		GithubToken: rOpts.GithubToken,
	}

	if rOpts.Plan != nil && rOpts.Plan.Release != nil {
		rlsOpts.ReleaseName = rOpts.Plan.Release.Name
	}
	if rOpts.Plan != nil && rOpts.Plan.Output != nil {
		rlsOpts.AssetFiles = []string{rOpts.Plan.Output.DiskImage}
	}

	return release.Check(c.Context, rlsOpts)
}
