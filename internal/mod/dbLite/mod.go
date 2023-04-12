package dbLite

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"supreme-flamego/core/kernel"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "dbLite"
}

func (m *Mod) Init(h *kernel.Hub) error {
	open, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return err
	}
	h.Map(&DB{open})
	return nil
}
