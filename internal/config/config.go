package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	mailerConfig "github.com/nash567/GoSentinel/internal/notifications/email/config"
	applicationConfig "github.com/nash567/GoSentinel/internal/service/application/config"
	authConfig "github.com/nash567/GoSentinel/internal/service/auth/config"
	cacheConfig "github.com/nash567/GoSentinel/pkg/cache/redis/config"
	dbConfig "github.com/nash567/GoSentinel/pkg/db/config"
	logConfig "github.com/nash567/GoSentinel/pkg/logger/config"
)

var ErrInvalidFileExtension = errors.New("file extension not supported")

type Config struct {
	HttpServer struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"http_server"`
	GrpcServer struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"grpc_server"`
	AppName           string                   `yaml:"app_name"`
	Env               string                   `yaml:"env"`
	DB                dbConfig.Config          `yaml:"database"`
	Log               logConfig.Config         `yaml:"logger"`
	Mailer            mailerConfig.Config      `yaml:"mailer"`
	ApplicationConfig applicationConfig.Config `yaml:"application_config"`
	AuthConfig        authConfig.Config        `yaml:"auth_config"`
	Redis             cacheConfig.Config       `yaml:"cache"`
}

func Load(fileNames ...string) (*Config, error) {
	loadFiles := make([]string, 0, len(fileNames))
	envFiles := make([]string, 0, len(fileNames))

	for _, file := range fileNames {
		fileParts := strings.Split(file, ".")
		fileExtn := fileParts[len(fileParts)-1]

		switch fileExtn {
		case "yml", "json", "yaml", "toml":
			loadFiles = append(loadFiles, file)
		case "env":
			envFiles = append(envFiles, file)
		default:
			return nil, ErrInvalidFileExtension
		}
	}

	if len(envFiles) > 0 {
		err := godotenv.Load(envFiles...)
		if err != nil {
			return nil, fmt.Errorf("error while loading env files(%s): %w", strings.Join(envFiles, ","), err)
		}
	}

	cfg, err := loadConfig(loadFiles...)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadConfig(fileNames ...string) (*Config, error) {
	var config Config

	err := configor.Load(&config, fileNames...)
	if err != nil {
		return nil, fmt.Errorf("cannot load config files(%s): %w", strings.Join(fileNames, ","), err)
	}

	return &config, nil
}
