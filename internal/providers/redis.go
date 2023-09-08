package providers

import (
	"context"
	"support/internal/env"

	"github.com/redis/go-redis/v9"
)

type DocumentDatabase interface {
	Connect() *redis.Client
	Set(key string, value interface{}) error
	Get(key string) (string, error)
}

type Redis struct {
	ctx    context.Context
	config env.EnvApp
}

func (r *Redis) Set(key string, value interface{}) error {
	rdb := r.Connect()
	err := rdb.Set(r.ctx, key, value, 0).Err()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	rdb := r.Connect()
	val, err := rdb.Get(r.ctx, key).Result()
	return val, err
}

func (r *Redis) Connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.config.REDIS_HOST + ":" + r.config.REDIS_PORT,
		Password: r.config.REDIS_PASSWORD,
		DB:       0,
	})
	return rdb
}

func NewRedis(ctx context.Context, ec env.EnvApp) DocumentDatabase {
	return &Redis{ctx: ctx, config: ec}
}
