package Logger

import (
	"log"

	"github.com/getsentry/sentry-go"
)

func Init() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://c804a58fb5d9385a47202b3a556a36c3@o4508080703864832.ingest.de.sentry.io/4508080888479824",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}
