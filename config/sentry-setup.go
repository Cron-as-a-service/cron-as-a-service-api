package config

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"os"
)

func InitSentry() {
	dsn := os.Getenv("SENTRY_DSN")
	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           dsn,
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		TracesSampler: sentry.TracesSampler(func(ctx sentry.SamplingContext) float64 {
			// As an example, this does not send some
			// transactions to Sentry based on their name.
			if ctx.Span.Name == "GET /healthz" {
				return 0.0
			}
			return 1.0
		}),
	}); err != nil {
		zap.L().Error(fmt.Sprintf("Sentry initialization failed: %v", err))
	}
}
