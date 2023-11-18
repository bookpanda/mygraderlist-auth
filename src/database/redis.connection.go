package database

import (
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func InitRedisConnect(conf *config.Redis) (cache *redis.Client, err error) {
	cache = redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		DB:       0,
		Password: conf.Password,
	})

	if cache == nil {
		return nil, errors.New("Cannot connect to redis server")
	}

	return
}
