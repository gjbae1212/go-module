package sentry

import (
	"errors"

	gosentry "github.com/getsentry/sentry-go"
)

var (
	ErrSentryEmptyParam = errors.New("[err] sentry empty params")
)

// InitSentry is to initialize Sentry setting.
func InitSentry(sentryDSN, environment, release, hostname string, stack, debug bool) error {
	if sentryDSN == "" || environment == "" || release == "" || hostname == "" {
		return ErrSentryEmptyParam
	}

	// if debug is true, it could show detail stack-log.
	if err := gosentry.Init(gosentry.ClientOptions{
		Dsn:              sentryDSN,
		Environment:      environment,
		Release:          release,
		ServerName:       hostname,
		AttachStacktrace: stack,
		Debug:            debug,
	}); err != nil {
		return err
	}

	return nil
}

// Error sends an error to Sentry.
func Error(err error) {
	if err == nil {
		return
	}
	gosentry.CaptureException(err)
}
