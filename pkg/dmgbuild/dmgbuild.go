package dmgbuild

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func Build(ctx context.Context, settings *Settings) error {
	if settings == nil {
		return fmt.Errorf("no settings provided")
	}

	logger := hclog.NewNullLogger()
	if settings.Logger != nil {
		logger = settings.Logger
	}

	if !strings.HasSuffix(logger.Name(), "dmgbuild") {
		logger = logger.Named("dmgbuild")
	}

	_, err := os.Stat(settings.Filename)
	if !os.IsNotExist(err) {
		return fmt.Errorf("output dmg exists: %s", settings.Filename)
	}

	baseCmd := settings.Command
	if baseCmd == "" {
		path, err2 := exec.LookPath("dmgbuild")
		if err2 != nil {
			return err2
		}
		baseCmd = path
	}

	file, err := settings.TempFile()
	if err != nil {
		return err
	}
	defer os.Remove(file)

	args := []string{"-s", file, settings.VolumeName, settings.Filename}

	if logger.IsDebug() {
		content, err2 := os.ReadFile(file)
		if err2 != nil {
			return err2
		}
		logger.Debug("using settings", file, string(content))
		logger.Debug("executing", "command", baseCmd, "args", args)
	}

	cmd := exec.CommandContext(ctx, baseCmd, args...)
	if settings.Stdout != nil {
		cmd.Stdout = settings.Stdout
	}
	if settings.Stderr != nil {
		cmd.Stderr = settings.Stderr
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	f, err := os.Stat(settings.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("output DMG file is missing")
		}

		return err
	}
	if !f.Mode().IsRegular() {
		return fmt.Errorf("output DMG file is not a file")
	}

	return nil
}
