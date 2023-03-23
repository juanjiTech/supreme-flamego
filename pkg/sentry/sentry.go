package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"supreme-flamego/conf"
)

func Init() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: conf.GetConfig().SentryDsn,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}
