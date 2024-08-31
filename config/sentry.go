package config

import (
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func InitSentry() {
	// Initialize Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
		EnableTracing:    true,
	})

	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)
}
