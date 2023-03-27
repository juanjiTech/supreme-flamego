package server

import (
	"fmt"
	"github.com/flamego/flamego"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"net"
	"os"
	"os/signal"
	"supreme-flamego/conf"
	"supreme-flamego/core/kernel"
	"supreme-flamego/core/logx"
	"supreme-flamego/internal/cache"
	"supreme-flamego/internal/database"
	"supreme-flamego/internal/mod/example"
	"supreme-flamego/internal/mod/flame"
	"supreme-flamego/pkg/colorful"
	"supreme-flamego/pkg/ip"
	"supreme-flamego/pkg/sentry"
	"syscall"
)

var log = logx.NameSpace("cmd.server")

var (
	configYml string
	e         *flamego.Flame
	StartCmd  = &cobra.Command{
		Use:     "server",
		Short:   "Set Application config info",
		Example: "main server -c ./config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("loading config...")
			conf.LoadConfig(configYml)
			log.Info("loading config complete")

			log.Info("init dep...")
			if conf.GetConfig().SentryDsn != "" {
				sentry.Init()
			}
			if conf.GetConfig().MODE == "" || conf.GetConfig().MODE == "debug" {
				logx.Init(zapcore.DebugLevel)
			} else {
				logx.Init(zapcore.InfoLevel)
			}
			database.InitDB()
			cache.InitCache()
			log.Info("init dep complete")

			log.Info("init kernel...")
			conn, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GetConfig().Port))
			if err != nil {
				log.Fatalw("failed to listen", "error", err)
			}
			tcpMux := cmux.New(conn)
			log.Infow("start listening", "port", conf.GetConfig().Port)
			k := kernel.New(kernel.Config{
				Listener: conn,
				MySQL:    database.GetDB("*")})
			k.Map(&tcpMux)
			k.RegMod(
				&example.Mod{},
				&flame.Mod{},
			)
			k.Init()
			log.Info("init kernel complete")

			log.Info("init module...")
			err = k.StartModule()
			if err != nil {
				panic(err)
			}
			log.Info("init module complete")

			log.Info("starting Server...")
			k.Serve()
			go func() {
				if err := tcpMux.Serve(); err != nil {
					panic(err)
				}
			}()

			fmt.Println(colorful.Green("Server run at:"))
			fmt.Println(fmt.Sprintf("-  Local:   http://localhost:%s", conf.GetConfig().Port))
			for _, host := range ip.GetLocalHost() {
				fmt.Println(fmt.Sprintf("-  Network: http://%s:%s", host, conf.GetConfig().Port))
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			fmt.Println(colorful.Blue("Shutting down server..."))

			err = k.Stop()
			if err != nil {
				panic(err)
			}
		},
	}
)

func init() {
	//StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/config.yaml", "Start server with provided configuration file")
}
