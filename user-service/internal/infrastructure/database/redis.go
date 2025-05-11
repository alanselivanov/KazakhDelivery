package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"user-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("Redis connected at %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return &RedisCache{
		Client: client,
		TTL:    time.Duration(cfg.Redis.TTL) * time.Second,
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, data, r.TTL).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *RedisCache) Close() error {
	log.Print("Closing Redis connection")
	return r.Client.Close()
}
