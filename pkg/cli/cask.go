package cli

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/jimeh/build-emacs-for-macos/pkg/cask"
	"github.com/jimeh/build-emacs-for-macos/pkg/repository"
	cli2 "github.com/urfave/cli/v2"
)

type caskOptions struct {
	BuildsRepo  *repository.Repository
	GithubToken string
}

func caskCmd() *cli2.Command {
	tokenDefaultText := ""
	if len(os.Getenv("GITHUB_TOKEN")) > 0 {
		tokenDefaultText = "***"
	}

	return &cli2.Command{
		Name:  "cask",
		Usage: "manage Homebrew Casks",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:    "builds-repository",
				Aliases: []string{"builds-repo", "b"},
				Usage:   "owner/name of GitHub repo for containing builds",
				EnvVars: []string{"EMACS_BUILDS_REPOSITORY"},
				Value:   "jimeh/emacs-builds",
			},
			&cli2.StringFlag{
				Name:        "github-token",
				Usage:       "GitHub API Token",
				EnvVars:     []string{"GITHUB_TOKEN"},
				DefaultText: tokenDefaultText,
				Required:    true,
			},
		},
		Subcommands: []*cli2.Command{
			caskUpdateCmd(),
		},
	}
}

func caskActionWrapper(
	f func(*cli2.Context, *Options, *caskOptions) error,
) func(*cli2.Context) error {
	return actionWrapper(func(c *cli2.Context, opts *Options) error {
		rOpts := &caskOptions{
			GithubToken: c.String("github-token"),
		}

		if r := c.String("builds-repository"); r != "" {
			var err error
			rOpts.BuildsRepo, err = repository.NewGitHub(r)
			if err != nil {
				return err
			}
		}

		return f(c, opts, rOpts)
	})
}

func caskUpdateCmd() *cli2.Command {
	return &cli2.Command{
		Name:      "update",
		Usage:     "update casks based on brew livecheck result in JSON format",
		ArgsUsage: "<livecheck.json>",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name: "ref",
				Usage: "git ref to create/update casks on top of in the " +
					"tap repository",
				EnvVars: []string{"GITHUB_REF"},
			},
			&cli2.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "directory to write cask files to",
			},
			&cli2.StringFlag{
				Name:    "tap-repository",
				Aliases: []string{"tap"},
				Usage: "owner/name of GitHub repo for Homebrew Tap to " +
					"commit changes to if --output is not set",
				EnvVars: []string{"GITHUB_REPOSITORY"},
			},
			&cli2.StringFlag{
				Name:     "templates-dir",
				Aliases:  []string{"t"},
				Usage:    "path to directory of cask templates",
				EnvVars:  []string{"CASK_TEMPLATE_DIR"},
				Required: true,
			},
			&cli2.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage: "force update file even if livecheck has it marked " +
					"as not outdated (does not force update if cask " +
					"content is unchanged)",
				Value: false,
			},
		},
		Action: caskActionWrapper(caskUpdateAction),
	}
}

func caskUpdateAction(
	c *cli2.Context,
	_ *Options,
	cOpts *caskOptions,
) error {
	updateOpts := &cask.UpdateOptions{
		BuildsRepo:   cOpts.BuildsRepo,
		GithubToken:  cOpts.GithubToken,
		Ref:          c.String("ref"),
		OutputDir:    c.String("output"),
		Force:        c.Bool("force"),
		TemplatesDir: c.String("templates-dir"),
	}

	if r := c.String("tap-repository"); r != "" {
		var err error
		updateOpts.TapRepo, err = repository.NewGitHub(r)
		if err != nil {
			return err
		}
	}

	arg := c.Args().First()
	if arg == "" {
		return errors.New("no livecheck argument given")
	}

	if arg == "-" {
		err := json.NewDecoder(c.App.Reader).Decode(&updateOpts.LiveChecks)
		if err != nil {
			return err
		}
	} else {
		f, err := os.Open(arg)
		if err != nil {
			return err
		}

		err = json.NewDecoder(f).Decode(&updateOpts.LiveChecks)
		if err != nil {
			return err
		}
	}

	return cask.Update(c.Context, updateOpts)
}
