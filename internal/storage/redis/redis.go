package redis

import (
	"context"
	"fmt"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Addr),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return rdb, nil

}
