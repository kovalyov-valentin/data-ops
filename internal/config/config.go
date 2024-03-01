package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	envPrefix   = ""
	envFilePath = "config/.env"
)

type Config struct {
	Env string
	HTTPServer
	PostgresDB
	Redis
	Nats
	Clickhouse
}

type HTTPServer struct {
	Address     string        `envconfig:"HTTP_SERVER_ADDRESS" default:"localhost:8080"`
	Timeout     time.Duration `envconfig:"HTTP_SERVER_TIMEOUT" default:"4s"`
	IdleTimeout time.Duration `envconfig:"HTTP_SERVER_CTX_TIMEOUT" default:"60s"`
	CtxTimeout  time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`
}

type PostgresDB struct {
	Username string `envconfig:"DB_USERNAME" default:"hezzl"`
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5040"`
	DBName   string `envconfig:"DB_NAME" default:"hezzldb"`
	Password string `envconfig:"DB_PASSWORD" default:"password"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

type Redis struct {
	Addr     string `envconfig:"REDIS_PORT" default:"REDIS_ADDR"`
	Password string `envconfig:"REDIS_PASSWORD" default:"password"`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
}

type Nats struct {
	Host string `envconfig:"NATS_HOST" default:"localhost"`
	Port string `envconfig:"NATS_PORT" default:"4222"`
}

type Clickhouse struct {
	//Addr       string `yaml:"addr" env:"CH_ADDR"`
	Username string `envconfig:"CLICKHOUSE_USERNAME" default:"hezzl"`
	Password string `envconfig:"CLICKHOUSE_PASSWORD" default:"password"`
	HttpPort int    `envconfig:"CLICKHOUSE_PORT" default:"8123"`
	DB       string `envconfig:"CLICKHOUSE_DB" default:"logs"`
	Host     string `envconfig:"CLICKHOUSE_HOST" default:"localhost"`
}

func MustLoad() *Config {
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err)
	}
	var cfg Config
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		logrus.Fatalf("Error filling in the structure: %s", err)
		return nil
	}
	return &cfg
}
