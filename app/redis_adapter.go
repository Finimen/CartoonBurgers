package main

import (
	"CartoonBurgers/services"
	"time"

	"github.com/go-redis/redis"
)

type RedisAdapter struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) services.IRedisClient {
	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Exists(key string) *redis.IntCmd {
	return r.client.Exists(key)
}

func (r *RedisAdapter) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *RedisAdapter) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(key, value, expiration)
}

func (r *RedisAdapter) Del(key string) *redis.IntCmd {
	return r.client.Del(key)
}
