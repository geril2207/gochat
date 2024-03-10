package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	ServerPort  string `envconfig:"SERVER_PORT"`
	DatabaseUrl string `envconfig:"DB_URL"`
	JwtKey      string `envconfig:"JWT_KEY"`
}

func NewConfigWithEnvPath(envPath string) (*EnvConfig, error) {
	_ = godotenv.Load(envPath)

	var config EnvConfig

	err := envconfig.Process("server", &config)

	return &config, err
}

func New() (*EnvConfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if strings.Contains(wd, "apps/server") {
		wd = filepath.Join(wd, "../..")
	}

	envPath := filepath.Join(wd, ".env")

	return NewConfigWithEnvPath(envPath)
}

func ProvideEnvConfig() EnvConfig {
	config, err := New()
	if err != nil {
		panic(err)
	}
	return *config
}
