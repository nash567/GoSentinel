package config

import "time"

type Config struct {
	Address              string `default:"localhost:6379"`
	Password             string
	DefaultKeyExpiryTime time.Duration `default:"1h"`
}
