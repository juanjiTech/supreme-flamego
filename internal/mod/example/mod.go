package example

import (
	"context"
	"fmt"
	"reflect"
	"supreme-flamego/core/kernel"
	"sync"
)

var _ kernel.Module = (*App)(nil)

type App struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (a *App) Name() string {
	return "example"
}

// 下面的方法皆为可选实现

func (a *App) PreInit(h *kernel.Hub) error {
	return nil
}

func (a *App) Init(h *kernel.Hub) error {
	h.Map("hello world") // 在内核注册这个依赖
	return nil
}

func (a *App) PostInit(h *kernel.Hub) error {
	return nil
}

func (a *App) Load(h *kernel.Hub) error {
	str := h.Value(reflect.TypeOf("string")).String() // 从内核获取上面注册的依赖
	fmt.Println(str)
	_, _ = h.Invoke(func(s string) { fmt.Println(s) }) // 也可以这样从内核获取上面注册的依赖

	h.Http.Get("/ping", func() string { return "pong" })
	return nil
}

func (a *App) Start(h *kernel.Hub) error {
	return nil
}

func (a *App) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	return nil
}
