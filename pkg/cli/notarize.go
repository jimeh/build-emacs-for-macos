package cli

import (
	"path/filepath"

	"github.com/jimeh/build-emacs-for-macos/pkg/notarize"
	"github.com/jimeh/build-emacs-for-macos/pkg/plan"
	cli2 "github.com/urfave/cli/v2"
)

func notarizeCmd() *cli2.Command {
	return &cli2.Command{
		Name:      "notarize",
		Usage:     "notarize and staple a dmg, zip, or pkg",
		ArgsUsage: "<file>",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:  "bundle-id",
				Usage: "bundle identifier",
				Value: "org.gnu.Emacs",
			},
			&cli2.StringFlag{
				Name:    "ac-username",
				Usage:   "Apple Connect username",
				EnvVars: []string{"AC_USERNAME"},
			},
			&cli2.StringFlag{
				Name:  "ac-password",
				Usage: "Apple Connect password",
				Value: "@env:AC_PASSWORD",
			},
			&cli2.StringFlag{
				Name:    "ac-provider",
				Usage:   "Apple Connect provider",
				EnvVars: []string{"AC_PROVIDER"},
			},
			&cli2.BoolFlag{
				Name:  "staple",
				Usage: "staple file after notarization",
				Value: true,
			},
			&cli2.StringFlag{
				Name: "plan",
				Usage: "path to build plan YAML file produced by " +
					"emacs-builder plan",
				Aliases:   []string{"p"},
				EnvVars:   []string{"EMACS_BUILDER_PLAN"},
				TakesFile: true,
			},
		},
		Action: actionWrapper(notarizeAction),
	}
}

func notarizeAction(c *cli2.Context, _ *Options) error {
	options := &notarize.Options{
		File:     c.Args().Get(0),
		BundleID: c.String("bundle-id"),
		Username: c.String("ac-username"),
		Password: c.String("ac-password"),
		Provider: c.String("ac-provider"),
		Staple:   c.Bool("staple"),
	}

	if f := c.String("plan"); f != "" {
		p, err := plan.Load(f)
		if err != nil {
			return err
		}

		if p.Output != nil {
			options.File = filepath.Join(
				p.Output.Directory, p.Output.DiskImage,
			)
		}
	}

	return notarize.Notarize(c.Context, options)
}
