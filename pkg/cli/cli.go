package cli

import (
	"fmt"
	"strings"

	cli2 "github.com/urfave/cli/v2"
)

type CLI struct {
	App     *cli2.App
	Version string
	Commit  string
	Date    string
}

func New(version, commit, date string) *CLI {
	if version == "" {
		version = "0.0.0-dev"
	}

	c := &CLI{
		Version: version,
		Commit:  commit,
		Date:    date,
		App: &cli2.App{
			Name:                 "emacs-builder",
			Usage:                "Tool to build emacs",
			Version:              version,
			EnableBashCompletion: true,
			Flags: []cli2.Flag{
				&cli2.StringFlag{
					Name:    "log-level",
					Usage:   "set log level",
					Aliases: []string{"l"},
					Value:   "info",
				},
				&cli2.BoolFlag{
					Name:    "quiet",
					Usage:   "silence noisy output",
					Aliases: []string{"q"},
					Value:   false,
				},
				cli2.VersionFlag,
			},
			Commands: []*cli2.Command{
				planCmd(),
				signCmd(),
				notarizeCmd(),
				packageCmd(),
				{
					Name:    "version",
					Usage:   "print the version",
					Aliases: []string{"v"},
					Action: func(c *cli2.Context) error {
						cli2.VersionPrinter(c)

						return nil
					},
				},
			},
		},
	}

	cli2.VersionPrinter = c.VersionPrinter

	return c
}

func (s *CLI) VersionPrinter(c *cli2.Context) {
	version := c.App.Version
	if version == "" {
		version = "0.0.0-dev"
	}

	extra := []string{}
	if len(s.Commit) >= 7 {
		extra = append(extra, s.Commit[0:7])
	}
	if s.Date != "" {
		extra = append(extra, s.Date)
	}
	var extraOut string
	if len(extra) > 0 {
		extraOut += " (" + strings.Join(extra, ", ") + ")"
	}

	fmt.Printf("%s%s\n", version, extraOut)
}

func (s *CLI) Run(args []string) error {
	return s.App.Run(args)
}
