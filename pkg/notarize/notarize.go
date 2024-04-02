package notarize

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/bearer/gon/notarize"
	"github.com/bearer/gon/staple"
	"github.com/hashicorp/go-hclog"
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
		File:        opts.File,
		DeveloperId: opts.Username,
		Password:    opts.Password,
		Provider:    opts.Provider,
		BaseCmd:     exec.CommandContext(ctx, ""),
		Status: &status{
			Lock:   &sync.Mutex{},
			Logger: logger,
		},
		// Ensure we don't log anything from the notarize package, as it will
		// log the password. We'll handle logging ourselves.
		Logger: hclog.NewNullLogger(),
	}

	logger.Info("notarizing", "file", filepath.Base(opts.File))

	info, _, err := notarize.Notarize(ctx, notarizeOpts)
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

	lastinfoStatus     string
	lastInfoStatusTime time.Time

	lastLogStatus     string
	lastLogStatusTime time.Time
}

var _ notarize.Status = (*status)(nil)

func (s *status) Submitting() {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Logger.Info("submitting file for notarization...")
}

func (s *status) Submitted(uuid string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	msg := "submitted, waiting for result from Apple"
	if s.Logger.IsDebug() {
		s.Logger.Debug(msg, "uuid", uuid)
	} else {
		s.Logger.Info(msg)
	}
}

func (s *status) InfoStatus(info notarize.Info) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.lastinfoStatus != info.Status ||
		time.Now().After(s.lastInfoStatusTime.Add(60*time.Second)) {
		s.lastinfoStatus = info.Status
		s.lastInfoStatusTime = time.Now()
		s.Logger.Info("status update", "status", info.Status)
	}
}

func (s *status) LogStatus(log notarize.Log) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.lastLogStatus != log.Status ||
		time.Now().After(s.lastLogStatusTime.Add(60*time.Second)) {
		s.lastLogStatus = log.Status
		s.lastLogStatusTime = time.Now()
		s.Logger.Info("log status update", "status", log.Status)
	}
}
