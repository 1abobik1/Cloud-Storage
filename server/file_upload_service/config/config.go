package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer        string        `env:"HTTP_SERVER_ADDRESS" env-required:"true"`
	JWTPublicKeyPath  string        `env:"JWT_PUBLIC_KEY_PATH" env-required:"true"`
	MinIoPort         string        `env:"MINIO_PORT" env-required:"true"`
	MinIoRootUser     string        `env:"MINIO_ROOT_USER" env-required:"true"`
	MinIoRootPassword string        `env:"MINIO_ROOT_PASSWORD" env-required:"true"`
	MinIoUseSSL       bool          `env:"MINIO_USE_SSL" env-required:"true"`
	MinIoURLLifeTime  time.Duration `env:"MINIO_URL_LIFETIME" env-required:"true"`
	RedisURLLifeTime  time.Duration `env:"REDIS_URL_LIFETIME" env-required:"true"`
	RedisPort         string        `env:"REDIS_PORT" env-required:"true"`
}

func MustLoad() *Config {
	path := getConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists" + path)
	}

	err := godotenv.Load(path)
	if err != nil {
		panic(fmt.Sprintf("No .env file found at %s, relying on environment variables", path))
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Sprintf("Failed to load environment variables: %v", err))
	}

	return &cfg
}

func getConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
