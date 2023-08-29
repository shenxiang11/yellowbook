package ioc

import (
	"github.com/redis/go-redis/v9"
	"yellowbook/config"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Conf.Redis.Addr,
	})
	return redisClient
}
