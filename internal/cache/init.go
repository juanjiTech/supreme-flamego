package cache

import (
	"supreme-flamego/conf"
	"supreme-flamego/core/logx"
	"supreme-flamego/internal/cache/types"
	"sync"
)

var (
	dbs = make(map[string]types.Cache)
	mux sync.RWMutex
)

func InitCache() {
	sources := conf.GetConfig().Caches
	for _, source := range sources {
		setCacheByKey(source.Key, mustCreateCache(source))
		if source.Key == "" {
			source.Key = "*"
		}
		logx.NameSpace("cache").Infof("create cache %s => %s:%s", source.Key, source.IP, source.PORT)
	}
}

func GetCache(key ...string) types.Cache {
	mux.Lock()
	defer mux.Unlock()
	if len(key) == 0 {
		return dbs["*"]
	}
	return dbs[key[0]]
}

func setCacheByKey(key string, cache types.Cache) {
	if key == "" {
		key = "*"
	}
	if GetCache(key) != nil {
		logx.NameSpace("cache").Error("duplicate db key: ", key)
	}
	mux.Lock()
	defer mux.Unlock()
	dbs[key] = cache
}

func mustCreateCache(conf conf.Cache) types.Cache {
	var creator = getCreatorByType(conf.Type)
	if creator == nil {
		logx.NameSpace("cache").Fatal("fail to find creator for cache types:%s", conf.Type)
		return nil
	}
	cache, err := creator.Create(conf)
	if err != nil {
		logx.NameSpace("cache").Fatal(err)
		return nil
	}
	return cache
}
