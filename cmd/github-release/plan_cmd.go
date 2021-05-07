package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var planCmd = &cli.Command{
	Name:      "plan",
	Usage:     "Plan if GitHub release and asset exists",
	UsageText: "github-release [global options] plan [<gif-ref>]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "emacs-mirror-repo",
			Usage:   "Github owner/repo to get Emacs commit info from",
			Aliases: []string{"e"},
			EnvVars: []string{"EMACS_MIRROR_REPO"},
			Value:   "emacs-mirror/emacs",
		},
	},
	Action: actionHandler(planAction),
}

func planAction(c *cli.Context, opts *globalOptions) error {
	gh := opts.gh
	planFile := opts.plan
	emacsRepo := NewRepo(c.String("emacs-mirror-repo"))

	ref := c.Args().Get(0)
	if ref == "" {
		ref = "master"
	}

	commit, _, err := gh.Repositories.GetCommit(
		c.Context, emacsRepo.Owner, emacsRepo.Name, ref,
	)
	if err != nil {
		return err
	}

	commitDate := commit.GetCommit().Committer.GetDate()
	date := commitDate.Format("2006-01-02")
	sha := commit.GetSHA()

	archive := fmt.Sprintf(
		"Emacs.app-[%s][%s][%s][%s][%s].tbz",
		nonAlphaNum.ReplaceAllString(ref, "-"),
		date,
		sha[0:7],
		"macOS-"+OSInfo.Version(),
		OSInfo.Arch(),
	)

	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	buildsDir := filepath.Join(rootDir, "builds")

	plan := &Plan{
		Ref:     ref,
		SHA:     sha,
		Date:    date,
		Archive: filepath.Join(buildsDir, archive),
	}

	b, err := yaml.Marshal(plan)
	if err != nil {
		return err
	}

	return os.WriteFile(planFile, b, 0666)
}
