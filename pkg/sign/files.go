package sign

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func Files(ctx context.Context, files []string, opts *Options) error {
	logger := hclog.FromContext(ctx).Named("sign")
	args := []string{}

	if opts.Identity != "" {
		args = append(args, "--sign", opts.Identity)
	}
	if opts.Deep {
		args = append(args, "--deep")
	}
	if opts.Timestamp {
		args = append(args, "--timestamp")
	}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}
	if len(opts.Options) > 0 {
		args = append(args, "--options", strings.Join(opts.Options, ","))
	}

	if opts.EntitlementsFile != "" {
		args = append(args, "--entitlements", opts.EntitlementsFile)
	} else if opts.Entitlements != nil {
		entitlementsFile, err := opts.Entitlements.TempFile()
		if err != nil {
			return err
		}
		defer os.Remove(entitlementsFile)
		logger.Debug("wrote entitlements", "file", entitlementsFile)

		args = append(args, "--entitlements", entitlementsFile)
	}

	baseCmd := opts.CodeSignCmd
	if baseCmd == "" {
		path, err := exec.LookPath("codesign")
		if err != nil {
			return err
		}
		baseCmd = path
	}

	args = append(args, files...)

	logger.Debug("executing", "command", baseCmd, "args", args)
	cmd := exec.CommandContext(ctx, baseCmd, args...)
	if opts.Output != nil {
		cmd.Stdout = opts.Output
		cmd.Stderr = opts.Output
	}

	return cmd.Run()
}
