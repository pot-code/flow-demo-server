package cache

import "github.com/redis/go-redis/v9"

func NewRedisCache(dsn string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
	})
	return rdb
}
