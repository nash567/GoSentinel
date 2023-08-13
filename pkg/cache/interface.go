package cache

import (
	"context"
	"time"
)

type KeyVal struct {
	Key    string
	Value  interface{}
	Expiry time.Duration
}

func NewKeyValWithExpiry(key string, value interface{}, expiry time.Duration) KeyVal {
	return KeyVal{
		Key:    key,
		Value:  value,
		Expiry: expiry,
	}
}

func NewKeyVal(key string, value interface{}) KeyVal {
	return KeyVal{
		Key:   key,
		Value: value,
	}
}

type Cache interface {
	Set(ctx context.Context, Key KeyVal) error
	Get(ctx context.Context, Key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetHealth() error
}
