package kernel

import (
	"context"
	"errors"
	"fmt"
	sentryflame "github.com/asjdf/flamego-sentry"
	"github.com/flamego/flamego"
	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"net/http"
	"supreme-flamego/conf"
	"supreme-flamego/core/discov"
	"supreme-flamego/core/logx"
	"supreme-flamego/pkg/colorful"
	"sync"
	"time"
)

type Engine struct {
	config Config
	// 下面这些将会转为private
	Mysql      *gorm.DB
	Cache      *redis.Client
	http       *flamego.Flame
	httpSrv    *http.Server
	grpc       *grpc.Server
	Conn       *grpc.ClientConn
	Mux        *runtime.ServeMux
	HttpServer *http.Server

	EtcdPublisher *discov.Publisher

	Ctx    context.Context
	Cancel context.CancelFunc

	ConfigListener []func(*conf.GlobalConfig)

	listener net.Listener

	modules   map[string]Module
	modulesMu sync.Mutex
}

type Config struct {
	Listener     net.Listener
	MySQL        *gorm.DB
	Redis        *redis.Client
	EnableSentry bool
}

func New(config ...Config) *Engine {
	if len(config) == 0 {
		panic("config can't be empty")
	}
	return &Engine{
		config:   config[0],
		listener: config[0].Listener,
		Mysql:    config[0].MySQL,
		Cache:    config[0].Redis,
		modules:  make(map[string]Module),
	}
}

func (e *Engine) Init() {
	e.Ctx, e.Cancel = context.WithCancel(context.Background())
	e.http = flamego.New()
	e.http.Use(flamego.Recovery())
	if e.config.EnableSentry {
		e.http.Use(sentryflame.New(sentryflame.Options{Repanic: true}))
	}
}

func (e *Engine) StartModule() error {
	hub := Hub{
		Http:  e.http,
		Grpc:  e.grpc,
		Mysql: e.Mysql,
		Cache: e.Cache,
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Logger = logx.NameSpace("module." + m.Name())
		if err := m.PreInit(&h4m); err != nil {
			return err
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Logger = logx.NameSpace("module." + m.Name())
		if err := m.Init(&h4m); err != nil {
			return err
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Logger = logx.NameSpace("module." + m.Name())
		if err := m.PostInit(&h4m); err != nil {
			return err
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Logger = logx.NameSpace("module." + m.Name())
		if err := m.Load(&h4m); err != nil {
			return err
		}
	}
	for _, m := range e.modules {
		h4m := hub
		h4m.Logger = logx.NameSpace("module." + m.Name())
		if err := m.Start(&h4m); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) Serve() {
	e.httpSrv = &http.Server{
		Handler: e.http,
	}

	go func() {
		if err := e.httpSrv.Serve(e.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(colorful.Red("Got Server Err: " + err.Error()))
		}
	}()
}

func (e *Engine) Stop() error {
	wg := sync.WaitGroup{}
	wg.Add(len(e.modules))
	for _, m := range e.modules {
		err := m.Stop(&wg, e.Ctx)
		if err != nil {
			return err
		}
	}
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.httpSrv.Shutdown(ctx); err != nil {
		fmt.Println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
		return err
	}

	return nil
}
