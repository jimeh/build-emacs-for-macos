package cli

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/dmg"
	"github.com/jimeh/build-emacs-for-macos/pkg/notarize"
	"github.com/jimeh/build-emacs-for-macos/pkg/plan"
	"github.com/jimeh/build-emacs-for-macos/pkg/sign"
	cli2 "github.com/urfave/cli/v2"
)

func packageCmd() *cli2.Command {
	return &cli2.Command{
		Name:      "package",
		Usage:     "package a build directory containing Emacs.app into a dmg",
		ArgsUsage: "<source-dir>",
		Flags: []cli2.Flag{
			&cli2.StringFlag{
				Name:    "volume-name",
				Usage:   "set volume name, defaults to basename of source dir",
				Aliases: []string{"n"},
			},
			&cli2.BoolFlag{
				Name: "sign",
				Usage: "sign Emacs.app before packaging, notarize and staple " +
					"dmg after packaging",
			},
			&cli2.StringFlag{
				Name: "output",
				Usage: "specify output dmg file name, if not specified the " +
					"output filename is based on source directory",
				Aliases: []string{"o"},
			},
			&cli2.BoolFlag{
				Name:    "sha256",
				Usage:   "create .sha256 checksum file for output dmg",
				Aliases: []string{"s"},
				Value:   true,
			},
			&cli2.BoolFlag{
				Name: "remove-source-dir",
				Usage: "remove source directory after successfully " +
					"creating dmg",
				Aliases: []string{"rm"},
				Value:   false,
			},
			&cli2.BoolFlag{
				Name:    "verbose",
				Usage:   "verbose output",
				Aliases: []string{"v"},
				Value:   false,
			},
			&cli2.StringFlag{
				Name:  "dmgbuild",
				Usage: "specify custom path to dmgbuild executable",
			},
			&cli2.StringFlag{
				Name:    "sign-identity",
				Usage:   "(with --sign) signing identity passed to codesign",
				EnvVars: []string{"AC_SIGN_IDENTITY"},
			},
			&cli2.StringFlag{
				Name:  "bundle-id",
				Usage: "(with --sign) bundle identifier",
				Value: "org.gnu.Emacs",
			},
			&cli2.StringFlag{
				Name:    "ac-username",
				Usage:   "(with --sign) Apple Connect username",
				EnvVars: []string{"AC_USERNAME"},
			},
			&cli2.StringFlag{
				Name:  "ac-password",
				Usage: "(with --sign) Apple Connect password",
				Value: "@env:AC_PASSWORD",
			},
			&cli2.StringFlag{
				Name:    "ac-provider",
				Usage:   "(with --sign) Apple Connect provider",
				EnvVars: []string{"AC_PROVIDER"},
			},
			&cli2.BoolFlag{
				Name:  "staple",
				Usage: "(with --sign) stable after notarization",
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
		Action: actionWrapper(packageAction),
	}
}

//nolint:funlen
func packageAction(c *cli2.Context, opts *Options) error {
	logger := hclog.FromContext(c.Context).Named("package")

	sourceDir := c.Args().Get(0)
	doSign := c.Bool("sign")

	var p *plan.Plan
	var err error
	if f := c.String("plan"); f != "" {
		p, err = plan.Load(f)
		if err != nil {
			return err
		}
	}

	if doSign {
		app := filepath.Join(sourceDir, "Emacs.app")

		signOpts := &sign.Options{
			Identity:  c.String("sign-identity"),
			Options:   []string{"runtime"},
			Deep:      true,
			Timestamp: true,
			Force:     true,
			Verbose:   c.Bool("verbose"),
		}

		if p != nil {
			if p.Output != nil && p.Build != nil {
				app = filepath.Join(
					p.Output.Directory, p.Build.Name, "Emacs.app",
				)
			}
		}

		if !opts.quiet {
			signOpts.Output = os.Stdout
		}

		err = sign.Emacs(c.Context, app, signOpts)
		if err != nil {
			return err
		}
	}

	dmgOpts := &dmg.Options{
		DMGBuild:        c.String("dmgbuild"),
		SourceDir:       sourceDir,
		VolumeName:      c.String("volume-name"),
		OutputFile:      c.String("output"),
		RemoveSourceDir: c.Bool("remove-source-dir"),
		Verbose:         c.Bool("verbose"),
	}

	if p != nil && p.Output != nil && p.Build != nil {
		dmgOpts.SourceDir = filepath.Join(
			p.Output.Directory, p.Build.Name,
		)
		dmgOpts.VolumeName = p.Build.Name
		dmgOpts.OutputFile = filepath.Join(
			p.Output.Directory, p.Output.DiskImage,
		)
	}

	if !opts.quiet {
		dmgOpts.Output = os.Stdout
	}

	outputDMG, err := dmg.Create(c.Context, dmgOpts)
	if err != nil {
		return err
	}

	if doSign {
		notarizeOpts := &notarize.Options{
			File:     outputDMG,
			BundleID: c.String("bundle-id"),
			Username: c.String("ac-username"),
			Password: c.String("ac-password"),
			Provider: c.String("ac-provider"),
			Staple:   c.Bool("staple"),
		}

		err = notarize.Notarize(c.Context, notarizeOpts)
		if err != nil {
			return err
		}
	}

	if c.Bool("sha256") {
		sumFile := outputDMG + ".sha256"

		logger.Info("generating SHA256 checksum", "file", outputDMG)
		sum, err := fileSHA256(outputDMG)
		if err != nil {
			return err
		}

		logger.Info("checksum", "sha256", sum, "file", outputDMG)
		content := fmt.Sprintf("%s  %s", sum, filepath.Base(outputDMG))
		err = os.WriteFile(sumFile, []byte(content), 0o644) //nolint:gosec
		if err != nil {
			return err
		}
		logger.Info("wrote checksum", "file", sumFile)
	}

	return nil
}

func fileSHA256(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
