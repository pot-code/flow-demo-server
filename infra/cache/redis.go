package cache

import "github.com/redis/go-redis/v9"

func NewRedisCache(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return rdb
}
