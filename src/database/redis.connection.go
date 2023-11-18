package database

import (
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func InitRedisConnect(conf *config.Redis) (cache *redis.Client, err error) {
	log.Info().Any("conf", conf)
	cache = redis.NewClient(&redis.Options{
		Addr:     conf.Host,
		DB:       0,
		Username: "",
		Password: conf.Password,
	})

	if cache == nil {
		return nil, errors.New("Cannot connect to redis server")
	}

	return
}
