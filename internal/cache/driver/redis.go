package driver

import (
	"fmt"
	"github.com/go-redis/redis"
	"supreme-flamego/conf"
	"supreme-flamego/core/logx"
	"supreme-flamego/internal/cache/types"
	"time"
)

type RedisCreator struct{}

func (c RedisCreator) Create(conf conf.Cache) (types.Cache, error) {
	var r RedisCache
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.IP, conf.PORT),
		Password: conf.PASSWORD,
		DB:       conf.DB,
	})
	_, err := r.client.Ping().Result()
	if err != nil {
		logx.NameSpace("redis").Fatal(err)
	}
	return r, nil
}

type RedisCache struct {
	client *redis.Client
}

func (r RedisCache) GetInt(key string) (int, bool) {
	value, err := r.client.Get(key).Int()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	return 0, false
}

func (r RedisCache) GetInt64(key string) (int64, bool) {
	value, err := r.client.Get(key).Int64()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	return 0, false
}

func (r RedisCache) GetFloat32(key string) (float32, bool) {
	value, err := r.client.Get(key).Float32()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	return 0, false
}

func (r RedisCache) GetFloat64(key string) (float64, bool) {
	value, err := r.client.Get(key).Float64()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	return 0, false
}

func (r RedisCache) GetString(key string) (string, bool) {
	value, err := r.client.Get(key).Result()
	if err == nil {
		return value, true
	}
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	return "", false
}

func (r RedisCache) GetBool(key string) (bool, bool) {
	value, err := r.client.Get(key).Result()
	if err != redis.Nil {
		logx.NameSpace("redis").Error(err)
	}
	if value == "1" {
		return true, true
	} else if value == "0" {
		return false, true
	}
	return false, false
}

func (r RedisCache) Set(Key string, value any, expireDuration time.Duration) error {
	return r.client.Set(Key, value, expireDuration).Err()
}

func (r RedisCache) Del(key string) bool {
	err := r.client.Del(key).Err()
	if err == redis.Nil {
		return false
	} else if err != nil {
		logx.NameSpace("redis").Error(err)
	}
	return true
}
