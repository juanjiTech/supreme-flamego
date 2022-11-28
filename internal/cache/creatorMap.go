package cache

import (
	"supreme-flamego/config"
	"supreme-flamego/internal/cache/driver"
	"supreme-flamego/internal/cache/types"
)

type Creator interface {
	Create(conf config.Cache) (types.Cache, error)
}

func init() {
	typeMap["redis"] = driver.RedisCreator{}
}

var typeMap = make(map[string]Creator)

func getCreatorByType(cacheType string) Creator {
	return typeMap[cacheType]
}
