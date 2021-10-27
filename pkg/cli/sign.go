package cli

import (
	"os"
	"path/filepath"

	"github.com/jimeh/build-emacs-for-macos/pkg/plan"
	"github.com/jimeh/build-emacs-for-macos/pkg/sign"
	cli2 "github.com/urfave/cli/v2"
)

func signCmd() *cli2.Command {
	return &cli2.Command{
		Name:      "sign",
		Usage:     "sign a Emacs.app bundle with codesign",
		ArgsUsage: "<emacs-app>",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:     "sign",
				Aliases:  []string{"s"},
				Usage:    "signing identity passed to codesign",
				EnvVars:  []string{"AC_SIGN_IDENTITY"},
				Required: true,
			},
			&cli2.StringSliceFlag{
				Name:    "entitlements",
				Aliases: []string{"e"},
				Usage:   "comma-separated list of entitlements to enable",
				Value:   cli2.NewStringSlice(sign.DefaultEmacsEntitlements...),
			},
			&cli2.BoolFlag{
				Name:    "deep",
				Aliases: []string{"d"},
				Usage:   "pass --deep to codesign",
				Value:   true,
			},
			&cli2.BoolFlag{
				Name:    "timestamp",
				Aliases: []string{"t"},
				Usage:   "pass --timestamp to codesign",
				Value:   true,
			},
			&cli2.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "pass --force to codesign",
				Value:   true,
			},
			&cli2.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "pass --verbose to codesign",
				Value:   false,
			},
			&cli2.StringSliceFlag{
				Name:    "options",
				Aliases: []string{"o"},
				Usage:   "options passed to codesign",
				Value:   cli2.NewStringSlice("runtime"),
			},
			&cli2.StringFlag{
				Name:  "codesign",
				Usage: "specify custom path to codesign executable",
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
		Action: actionWrapper(signAction),
	}
}

func signAction(c *cli2.Context, opts *Options) error {
	signOpts := &sign.Options{
		Identity:    c.String("sign"),
		Options:     c.StringSlice("options"),
		Deep:        c.Bool("deep"),
		Timestamp:   c.Bool("timestamp"),
		Force:       c.Bool("force"),
		Verbose:     c.Bool("verbose"),
		CodeSignCmd: c.String("codesign"),
	}

	if v := c.StringSlice("entitlements"); len(v) > 0 {
		e := sign.Entitlements(v)
		signOpts.Entitlements = &e
	}

	if !opts.quiet {
		signOpts.Output = os.Stdout
	}

	app := c.Args().Get(0)

	if f := c.String("plan"); f != "" {
		p, err := plan.Load(f)
		if err != nil {
			return err
		}

		if p.Output != nil && p.Build != nil {
			app = filepath.Join(
				p.Output.Directory, p.Build.Name, "Emacs.app",
			)
		}
	}

	return sign.Emacs(c.Context, app, signOpts)
}

func signFilesCmd() *cli2.Command {
	signCmd := signCmd()

	var flags []cli2.Flag
	for _, f := range signCmd.Flags {
		n := f.Names()
		if len(n) > 0 && n[0] == "plan" {
			continue
		}

		flags = append(flags, f)
	}

	return &cli2.Command{
		Name:      "sign-files",
		Usage:     "sign files with codesign",
		ArgsUsage: "<file> [<file>...]",
		Hidden:    true,
		Flags:     flags,
		Action:    actionWrapper(signFilesAction),
	}
}

func signFilesAction(c *cli2.Context, opts *Options) error {
	signOpts := &sign.Options{
		Identity:    c.String("sign"),
		Options:     c.StringSlice("options"),
		Deep:        c.Bool("deep"),
		Timestamp:   c.Bool("timestamp"),
		Force:       c.Bool("force"),
		Verbose:     c.Bool("verbose"),
		CodeSignCmd: c.String("codesign"),
	}

	if v := c.StringSlice("entitlements"); len(v) > 0 {
		e := sign.Entitlements(v)
		signOpts.Entitlements = &e
	}

	if !opts.quiet {
		signOpts.Output = os.Stdout
	}

	return sign.Files(c.Context, c.Args().Slice(), signOpts)
}
