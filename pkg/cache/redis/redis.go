package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nash567/GoSentinel/pkg/cache"

	"github.com/go-redis/redis/v8"
	"github.com/nash567/GoSentinel/pkg/cache/redis/config"
)

const redisExpirationTimeInSecond = 5

type Redis struct {
	client               *redis.Client
	defaultKeyExpiryTime time.Duration
}

func New(ctx context.Context, cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection: %w", err)
	}

	return &Redis{
		client:               client,
		defaultKeyExpiryTime: cfg.DefaultKeyExpiryTime,
	}, nil
}

func (r Redis) Set(ctx context.Context, keyVal cache.KeyVal) error {
	resp, err := json.Marshal(keyVal.Value)
	if err != nil {
		return fmt.Errorf("marshal file status: %w", err)
	}
	expiry := keyVal.Expiry
	if expiry == 0 {
		expiry = r.defaultKeyExpiryTime
	}
	err = r.client.Set(ctx, keyVal.Key, resp, expiry).Err()
	if err != nil {
		return fmt.Errorf("set value in redis: %w", err)
	}
	return nil
}

func (r Redis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("get value from redis: %w", err)
	}

	return value, nil
}

func (r Redis) GetHealth() error {
	err := r.client.Set(context.TODO(), "health-test", "value", time.Duration(redisExpirationTimeInSecond)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("redis health test: %w", err)
	}

	err = r.client.Del(context.TODO(), "health-test").Err()
	if err != nil {
		return fmt.Errorf("redis health test: %w", err)
	}
	return nil
}

func (r Redis) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("delete value from redis: %w", err)
	}

	return nil
}
