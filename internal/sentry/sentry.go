package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"supreme-flamego/config"
)

func Init() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetConfig().SentryDsn,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}
