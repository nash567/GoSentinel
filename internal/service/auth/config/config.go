package config

import "time"

type Config struct {
	JwtSecret     string
	JWTExpiration time.Duration
	EncryptionKey string
	SecretLength  int
}
