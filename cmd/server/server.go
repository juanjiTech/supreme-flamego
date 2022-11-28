package server

import (
	"context"
	"errors"
	"fmt"
	sentryflame "github.com/asjdf/flamego-sentry"
	"github.com/flamego/flamego"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"supreme-flamego/config"
	"supreme-flamego/internal/app/routerInitialize"
	"supreme-flamego/internal/cache"
	"supreme-flamego/internal/database"
	"supreme-flamego/internal/sentry"
	"supreme-flamego/pkg/colorful"
	"supreme-flamego/pkg/ip"
	"supreme-flamego/pkg/logger"
	"syscall"
	"time"
)

var (
	configYml string
	e         *flamego.Flame
	StartCmd  = &cobra.Command{
		Use:     "server",
		Short:   "Set Application config info",
		Example: "main server -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			println("loading config...")
			setUp()
			println("loading config complete")
			println("loading api...")
			load()
			println("loading api complete")
		},
		Run: func(cmd *cobra.Command, args []string) {
			println("starting Server...")
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}

func setUp() {
	// 顺序不能变 logger依赖config logger后面的同时依赖logger和config 否则crash
	config.LoadConfig(configYml)

	if config.GetConfig().SentryDsn != "" {
		sentry.Init()
	}

	if config.GetConfig().MODE == "" || config.GetConfig().MODE == "debug" {
		logger.Init(zapcore.DebugLevel)
	} else {
		logger.Init(zapcore.InfoLevel)
	}

	database.InitDB()
	cache.InitCache()
}

func load() {
	flamego.SetEnv(flamego.EnvType(config.GetConfig().MODE))
	e = flamego.New()
	e.Use(flamego.Recovery())
	if config.GetConfig().SentryDsn != "" {
		e.Use(sentryflame.New(sentryflame.Options{Repanic: true}))
	}

	routerInitialize.ApiInit(e)
}

func run() {
	srv := &http.Server{
		Addr:    config.GetConfig().Listen + ":" + config.GetConfig().Port,
		Handler: e,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			println(colorful.Red("Got Server Err: " + err.Error()))
		}
	}()

	println(colorful.Green("Server run at:"))
	println(fmt.Sprintf("-  Local:   http://localhost:%s", config.GetConfig().Port))
	for _, host := range ip.GetLocalHost() {
		println(fmt.Sprintf("-  Network: http://%s:%s", host, config.GetConfig().Port))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	println(colorful.Blue("Shutting down server..."))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
	}

	println(colorful.Green("Server exiting Correctly"))
}
