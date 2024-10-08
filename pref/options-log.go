package pref

import (
	"log/slog"
)

type (
	MonitorOptions struct {
		Log *slog.Logger
	}

	LogRotationOptions struct {
		// MaxSizeInMb, max size of a log file, before it is re-cycled
		MaxSizeInMb int

		// MaxNoOfBackups, max number of legacy log files that can exist
		// before being deleted
		MaxNoOfBackups int

		// MaxAgeInDays, max no of days before old log file is deleted
		MaxAgeInDays int
	}
)

// WithLogger defines a structure logger
func WithLogger(logger *slog.Logger) Option {
	return func(o *Options) error {
		o.Monitor.Log = logger

		return nil
	}
}
