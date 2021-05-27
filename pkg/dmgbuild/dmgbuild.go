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
	logger := hclog.NewNullLogger()
	if settings.Logger == nil {
		logger = settings.Logger
	}

	if !strings.HasSuffix(logger.Name(), "dmgbuild") {
		logger = logger.Named("dmgbuild")
	}

	if settings == nil {
		return fmt.Errorf("no settings provided")
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
		content, err := os.ReadFile(file)
		if err != nil {
			return err
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

	return cmd.Run()
}
