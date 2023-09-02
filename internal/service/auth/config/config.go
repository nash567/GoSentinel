package config

import "time"

type Config struct {
	JwtSecret                 string
	VerificationJWTExpiration time.Duration
	EncryptionKey             string
	SecretLength              int
	ApplicationJWTExpiry      time.Duration
	UserJWTExpiry             time.Duration
}
