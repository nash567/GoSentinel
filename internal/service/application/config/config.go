package config

import "time"

type Config struct {
	VerificationTemplate string
	VerificationExpiry   time.Duration
}
