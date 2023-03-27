package example

import (
	"context"
	"errors"
	"fmt"
	"github.com/flamego/flamego"
	"reflect"
	"supreme-flamego/core/kernel"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "example"
}

// 下面的方法皆为可选实现

func (m *Mod) PreInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Init(h *kernel.Hub) error {
	h.Map("hello world") // 在内核注册这个依赖
	return nil
}

func (m *Mod) PostInit(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Load(h *kernel.Hub) error {
	str := h.Value(reflect.TypeOf("string")).String() // 从内核获取上面注册的依赖
	fmt.Println(str)
	_, _ = h.Invoke(func(s string) { fmt.Println(s) }) // 也可以这样从内核获取上面注册的依赖

	var http *flamego.Flame
	http, ok := h.Value(reflect.TypeOf(http)).Interface().(*flamego.Flame)
	if !ok {
		return errors.New("flameGo didn't injected to kernel")
	}
	http.Get("/ping", func() string { return "pong" })
	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	return nil
}
