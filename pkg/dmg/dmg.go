package dmg

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/jimeh/build-emacs-for-macos/pkg/dmg/assets"
	"github.com/jimeh/build-emacs-for-macos/pkg/dmgbuild"
)

type Options struct {
	DMGBuild string

	SourceDir       string
	VolumeName      string
	OutputFile      string
	RemoveSourceDir bool
	Verbose         bool
	Output          io.Writer
}

// Create will create a *.dmg disk image as specified by the given Options.
//
//nolint:funlen
func Create(ctx context.Context, opts *Options) (string, error) {
	logger := hclog.FromContext(ctx).Named("package")

	sourceDir, err := filepath.Abs(opts.SourceDir)
	if err != nil {
		return "", err
	}

	appBundle := filepath.Join(sourceDir, "Emacs.app")
	_, err = os.Stat(appBundle)
	if err != nil {
		return "", err
	}

	volIcon, err := assets.IconTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(volIcon)

	bgImg, err := assets.BackgroundTempFile()
	if err != nil {
		return "", err
	}
	defer os.Remove(bgImg)

	volName := opts.VolumeName
	if volName == "" {
		volName = filepath.Base(sourceDir)
	}

	outputDMG := opts.OutputFile
	if outputDMG == "" {
		outputDMG = sourceDir + ".dmg"
	}

	settings := &dmgbuild.Settings{
		Logger: logger,

		Filename:         outputDMG,
		VolumeName:       volName,
		Icon:             volIcon,
		Format:           dmgbuild.UDZOFormat,
		CompressionLevel: 9,
		Files: []*dmgbuild.File{
			{
				Path: appBundle,
				PosX: 170,
				PosY: 200,
			},
		},
		Symlinks: []*dmgbuild.Symlink{
			{
				Name:   "Applications",
				Target: "/Applications",
				PosX:   510,
				PosY:   200,
			},
		},
		Window: dmgbuild.Window{
			Background:  bgImg,
			PoxX:        200,
			PosY:        200,
			Width:       680,
			Height:      446,
			DefaultView: dmgbuild.Icon,
		},
		IconView: dmgbuild.IconView{
			IconSize: 160,
			TextSize: 16,
		},
	}

	copyingFile := filepath.Join(sourceDir, "COPYING")
	fi, err := os.Stat(copyingFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	} else if err == nil && fi.Mode().IsRegular() {
		settings.Files = append(settings.Files, &dmgbuild.File{
			Path: copyingFile,
			PosX: 340,
			PosY: 506,
		})
	}

	if opts.Output != nil {
		settings.Stdout = opts.Output
		settings.Stderr = opts.Output
	}

	logger.Info("creating dmg", "file", filepath.Base(outputDMG))

	err = dmgbuild.Build(ctx, settings)
	if err != nil {
		return "", err
	}

	if opts.RemoveSourceDir {
		dir, err := filepath.Abs(opts.SourceDir)
		if err != nil {
			return "", err
		}

		logger.Info("removing", "source-dir", dir)
		err = os.RemoveAll(dir)
		if err != nil {
			return "", err
		}
	}

	return outputDMG, nil
}
