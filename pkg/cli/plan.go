package cli

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/plan"
	cli2 "github.com/urfave/cli/v2"
)

func planCmd() *cli2.Command {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	tokenDefaultText := ""
	if len(os.Getenv("GITHUB_TOKEN")) > 0 {
		tokenDefaultText = "***"
	}

	return &cli2.Command{
		Name:      "plan",
		Usage:     "plan a Emacs.app bundle with codeplan",
		ArgsUsage: "<branch/tag>",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name: "emacs-repo",
				Usage: "GitHub repository to get Emacs commit info and " +
					"tarball from",
				Aliases: []string{"e"},
				EnvVars: []string{"EMACS_REPO"},
				Value:   "emacs-mirror/emacs",
			},
			&cli2.StringFlag{
				Name:  "sha",
				Usage: "override commit SHA of specified git branch/tag",
			},
			&cli2.StringFlag{
				Name: "output",
				Usage: "output filename to write plan to instead of printing " +
					"to STDOUT",
				Aliases: []string{"o"},
			},
			&cli2.StringFlag{
				Name:  "output-dir",
				Usage: "output directory where build result is stored",
				Value: filepath.Join(wd, "builds"),
			},
			&cli2.StringFlag{
				Name: "test-build",
				Usage: "plan a test build with given name, which is " +
					"published to a draft or pre-release " +
					"\"test-builds\" release",
			},
			&cli2.StringFlag{
				Name:  "test-release-type",
				Value: "prerelease",
				Usage: "type of release when doing a test-build " +
					"(prerelease or draft)",
			},
			&cli2.StringFlag{
				Name:        "github-token",
				Usage:       "GitHub API Token",
				EnvVars:     []string{"GITHUB_TOKEN"},
				DefaultText: tokenDefaultText,
			},
		},
		Action: actionWrapper(planAction),
	}
}

func planAction(c *cli2.Context, opts *Options) error {
	logger := hclog.FromContext(c.Context).Named("plan")

	ref := c.Args().Get(0)
	if ref == "" {
		ref = "master"
	}

	planOpts := &plan.Options{
		EmacsRepo:     c.String("emacs-repo"),
		Ref:           ref,
		SHAOverride:   c.String("sha"),
		OutputDir:     c.String("output-dir"),
		TestBuild:     c.String("test-build"),
		TestBuildType: plan.Prerelease,
		GithubToken:   c.String("github-token"),
	}

	if c.String("test-release-type") == "draft" {
		planOpts.TestBuildType = plan.Draft
	}

	if !opts.quiet {
		planOpts.Output = os.Stdout
	}

	p, err := plan.Create(c.Context, planOpts)
	if err != nil {
		return err
	}

	planYAML, err := p.YAML()
	if err != nil {
		return err
	}

	var out *os.File
	out = os.Stdout
	if f := c.String("output"); f != "" {
		logger.Info("writing plan", "file", f)
		logger.Debug("content", "yaml", planYAML)
		out, err = os.Create(f)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	_, err = out.WriteString(planYAML)
	if err != nil {
		return err
	}

	return nil
}
