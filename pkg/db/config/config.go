package config

import "time"

type Config struct {
	UserName        string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host" env:"DB_HOST"`
	Port            int           `yaml:"port" env:"DB_PORT"`
	DBName          string        `yaml:"name" env:"DB_NAME"`
	MaxConnLifeTime time.Duration `yaml:"max-conn-life" env:"DB_CONN_MAX_LIFE_TIME"`
	MaxConns        int           `yaml:"max-conns" env:"DB_MAX_CONNS"`
	MaxIdleConns    int           `yaml:"max-idle-conns" env:"DB_MAX_IDLE_CONNS"`
}
