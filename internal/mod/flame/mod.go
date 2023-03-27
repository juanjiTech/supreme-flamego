package flame

import (
	"context"
	"github.com/flamego/flamego"
	"supreme-flamego/core/kernel"
	"sync"
)

var _ kernel.Module = (*App)(nil)

type App struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (a *App) Name() string {
	return "flame"
}

// 下面的方法皆为可选实现

func (a *App) PreInit(h *kernel.Hub) error {
	return nil
}

func (a *App) Init(h *kernel.Hub) error {
	h.Map(flamego.New()) // 在内核注册这个依赖
	return nil
}

func (a *App) PostInit(h *kernel.Hub) error {
	return nil
}

func (a *App) Load(h *kernel.Hub) error {
	return nil
}

func (a *App) Start(h *kernel.Hub) error {
	return nil
}

func (a *App) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	return nil
}
