package config

import (
	"fmt"
	"goalify/utils/options"
	"os"

	"github.com/joho/godotenv"
)

type ConfigService struct {
	envFilePath string
}

func NewConfigService(envFilePath options.Option[string]) *ConfigService {
	if os.Getenv("CI") == "true" {
		return &ConfigService{envFilePath: ""}
	}
	var err error
	if envFilePath.IsPresent() {
		err = godotenv.Load(envFilePath.ValueOrZero())
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		panic(err)
	}

	return &ConfigService{envFilePath: envFilePath.ValueOrZero()}
}

func (c *ConfigService) GetEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("env variable %s not found", key)
	}
	return val, nil
}

func (c *ConfigService) MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("env variable %s not found", key))
	}
	return val
}

func (c *ConfigService) SetEnv(key, value string) {
	os.Setenv(key, value)
}
