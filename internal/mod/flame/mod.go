package flame

import (
	"context"
	"errors"
	"fmt"
	sentryflame "github.com/asjdf/flamego-sentry"
	"github.com/charmbracelet/log"
	"github.com/flamego/flamego"
	"github.com/juanjiTech/inject"
	"github.com/soheilhy/cmux"
	"net/http"
	"reflect"
	"supreme-flamego/conf"
	"supreme-flamego/core/kernel"
	"supreme-flamego/pkg/colorful"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule

	flame   *flamego.Flame
	httpSrv *http.Server
}

func (m *Mod) Name() string {
	return "flame"
}

// 下面的方法皆为可选实现

func (m *Mod) PreInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Init(h *kernel.Hub) error {
	m.flame = flamego.New()
	m.flame.Use(flamego.LoggerInvoker(func(c flamego.Context, log *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %s", err)

				// Lookup the current ResponseWriter
				val := c.Value(inject.InterfaceOf((*http.ResponseWriter)(nil)))
				w := val.Interface().(http.ResponseWriter)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}))
	if conf.GetConfig().SentryDsn != "" {
		m.flame.Use(sentryflame.New(sentryflame.Options{Repanic: true}))
	}
	h.Map(m.flame) // 在内核注册这个依赖
	return nil
}

func (m *Mod) PostInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Load(h *kernel.Hub) error {
	var flame *flamego.Flame
	flame, ok := h.Value(reflect.TypeOf(flame)).Interface().(*flamego.Flame)
	if !ok {
		return errors.New("flameGo didn't injected to kernel")
	}
	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	var tcpMux *cmux.CMux
	tcpMux, ok := h.Value(reflect.TypeOf(tcpMux)).Interface().(*cmux.CMux)
	if !ok {
		return errors.New("flameGo didn't injected to kernel")
	}

	httpL := (*tcpMux).Match(cmux.HTTP1Fast())
	m.httpSrv = &http.Server{
		Handler: m.flame,
	}

	if err := m.httpSrv.Serve(httpL); err != nil && err != http.ErrServerClosed {
		h.Logger.Infow("failed to start to listen and serve", "error", err)
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.httpSrv.Shutdown(ctx); err != nil {
		fmt.Println(colorful.Yellow("Server forced to shutdown: " + err.Error()))
		return err
	}
	return nil
}
