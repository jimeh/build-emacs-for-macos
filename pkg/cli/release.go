package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
			releasePublishCmd(),
			releaseBulkCmd(),
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

func releasePublishCmd() *cli2.Command {
	return &cli2.Command{
		Name: "publish",
		Usage: "publish a GitHub release with specified asset " +
			"files",
		ArgsUsage: "[<asset-file> ...]",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:    "sha",
				Aliases: []string{"s"},
				Usage:   "git SHA to create release on",
				EnvVars: []string{"GITHUB_SHA"},
			},
			&cli2.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "release type, must be normal, prerelease, or draft",
				Value:   "normal",
			},
			&cli2.StringFlag{
				Name: "title",
				Usage: "release title, will use release name if not " +
					"specified",
				Value: "",
			},
		},
		Action: releaseActionWrapper(releasePublishAction),
	}
}

func releasePublishAction(
	c *cli2.Context,
	opts *Options,
	rOpts *releaseOptions,
) error {
	rlsOpts := &release.PublishOptions{
		Repository:   rOpts.Repository,
		CommitRef:    c.String("release-sha"),
		ReleaseName:  rOpts.Name,
		ReleaseTitle: c.String("title"),
		AssetFiles:   c.Args().Slice(),
		GithubToken:  rOpts.GithubToken,
	}

	rlsType := c.String("type")
	switch rlsType {
	case "draft":
		rlsOpts.ReleaseType = release.Draft
	case "prerelease":
		rlsOpts.ReleaseType = release.Prerelease
	case "normal":
		rlsOpts.ReleaseType = release.Normal
	default:
		return fmt.Errorf("invalid --type \"%s\"", rlsType)
	}

	if rOpts.Plan != nil {
		if rOpts.Plan.Release != nil {
			rlsOpts.ReleaseName = rOpts.Plan.Release.Name
			rlsOpts.ReleaseTitle = rOpts.Plan.Release.Title

			if rOpts.Plan.Release.Draft {
				rlsOpts.ReleaseType = release.Draft
			} else if rOpts.Plan.Release.Prerelease {
				rlsOpts.ReleaseType = release.Prerelease
			}
		}

		if rOpts.Plan.Output != nil {
			rlsOpts.AssetFiles = []string{
				filepath.Join(
					rOpts.Plan.Output.Directory,
					rOpts.Plan.Output.DiskImage,
				),
			}
		}
	}

	return release.Publish(c.Context, rlsOpts)
}

func releaseBulkCmd() *cli2.Command {
	return &cli2.Command{
		Name:      "bulk",
		Usage:     "bulk modify GitHub releases",
		ArgsUsage: "",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:  "name",
				Usage: "regexp pattern matching release names to modify",
			},
			&cli2.StringFlag{
				Name: "prerelease",
				Usage: "change prerelease flag, must be \"true\" or " +
					"\"false\", otherwise prerelease value is not changed",
			},
			&cli2.BoolFlag{
				Name:  "dry-run",
				Usage: "do not perform any changes",
			},
		},
		Action: releaseActionWrapper(releaseBulkAction),
	}
}

func releaseBulkAction(
	c *cli2.Context,
	opts *Options,
	rOpts *releaseOptions,
) error {
	bulkOpts := &release.BulkOptions{
		Repository:  rOpts.Repository,
		NamePattern: c.String("name"),
		DryRun:      c.Bool("dry-run"),
		GithubToken: rOpts.GithubToken,
	}

	switch c.String("prerelease") {
	case "true":
		v := true
		bulkOpts.Prerelease = &v
	case "false":
		v := false
		bulkOpts.Prerelease = &v
	case "":
	default:
		return errors.New(
			"--prerelease by me \"true\" or \"false\" when specified",
		)
	}

	return release.Bulk(c.Context, bulkOpts)
}
