package cli

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	cli2 "github.com/urfave/cli/v2"
)

type Options struct {
	quiet bool
}

func actionWrapper(
	f func(*cli2.Context, *Options) error,
) func(*cli2.Context) error {
	return func(c *cli2.Context) error {
		opts := &Options{
			quiet: c.Bool("quiet"),
		}

		levelStr := c.String("log-level")
		level := hclog.LevelFromString(levelStr)
		if level == hclog.NoLevel {
			return fmt.Errorf("invalid log level \"%s\"", levelStr)
		}

		// Prevent things from logging if they weren't explicitly given a
		// logger.
		hclog.SetDefault(hclog.NewNullLogger())

		// Create custom logger.
		logr := hclog.New(&hclog.LoggerOptions{
			Level:      level,
			Output:     os.Stderr,
			Mutex:      &sync.Mutex{},
			TimeFormat: time.RFC3339,
			Color:      hclog.ColorOff,
		})

		ctx := hclog.WithContext(c.Context, logr)
		c.Context = ctx

		return f(c, opts)
	}
}
