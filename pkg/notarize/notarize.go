package notarize

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/gon/notarize"
	"github.com/mitchellh/gon/staple"
)

type Options struct {
	File     string
	BundleID string
	Username string
	Password string
	Provider string
	Staple   bool
}

func Notarize(ctx context.Context, opts *Options) error {
	logger := hclog.FromContext(ctx).Named("notarize")

	notarizeOpts := &notarize.Options{
		File:     opts.File,
		BundleId: opts.BundleID,
		Username: opts.Username,
		Password: opts.Password,
		Provider: opts.Provider,
		BaseCmd:  exec.CommandContext(ctx, ""),
		Status: &status{
			Lock:   &sync.Mutex{},
			Logger: logger,
		},
	}

	logger.Info("notarizing", "file", filepath.Base(opts.File))

	info, err := notarize.Notarize(ctx, notarizeOpts)
	if err != nil {
		return err
	}

	logger.Info(
		"notarization complete",
		"status", info.Status,
		"message", info.StatusMessage,
	)

	if opts.Staple {
		logger.Info("stapling", "file", filepath.Base(opts.File))
		err := staple.Staple(ctx, &staple.Options{
			File:    opts.File,
			BaseCmd: exec.CommandContext(ctx, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type status struct {
	Lock   *sync.Mutex
	Logger hclog.Logger

	lastStatusTime time.Time
}

func (s *status) Submitting() {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Logger.Info("submitting file for notarization...")
}

func (s *status) Submitted(uuid string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Logger.Info("submitted")
	s.Logger.Debug("request", "uuid", uuid)
	s.Logger.Info("waiting for result from Apple...")
}

func (s *status) Status(info notarize.Info) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if time.Now().After(s.lastStatusTime.Add(60 * time.Second)) {
		s.lastStatusTime = time.Now()
		s.Logger.Info("status update", "status", info.Status)
	}
}
