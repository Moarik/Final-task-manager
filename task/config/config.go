package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"taskManager/task/pkg/postgre"
	"time"
)

type (
	Config struct {
		Postgres postgre.Config
		Server   Server
		Nats     Nats
		Cache    Cache
	}

	Cache struct {
		ClientTTL time.Duration `env:"REDIS_CACHE_CLIENT_TTL" envDefault:"24h"`

		CMSVariableRefreshTime time.Duration `env:"CLIENT_REFRESH_TIME" envDefault:"1m"`
	}

	Server struct {
		HTTPServer HTTPServer
		GRPCServer GRPCServer
	}

	HTTPServer struct {
		Port           int           `env:"HTTP_PORT,required"`
		ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"30s"`
		WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
		IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
		MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES" envDefault:"1048576"` // 1 MB
		TrustedProxies []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
		Mode           string        `env:"GIN_MODE" envDefault:"release"` // Can be: release, debug, test
	}

	GRPCServer struct {
		Port                  int16         `env:"GRPC_PORT"` // removed notEmpty
		MaxRecvMsgSizeMiB     int           `env:"GRPC_MAX_MESSAGE_SIZE_MIB" envDefault:"12"`
		MaxConnectionAge      time.Duration `env:"GRPC_MAX_CONNECTION_AGE" envDefault:"30s"`
		MaxConnectionAgeGrace time.Duration `env:"GRPC_MAX_CONNECTION_AGE_GRACE" envDefault:"10s"`
	}

	Nats struct {
		Hosts        []string `env:"NATS_HOSTS,notEmpty" envSeparator:","`
		NKey         string   `env:"NATS_NKEY"`
		IsTest       bool     `env:"NATS_IS_TEST,notEmpty" envDefault:"true"`
		NatsSubjects NatsSubjects
	}

	NatsSubjects struct {
		ClientEventSubject string `env:"NATS_CLIENT_EVENT_SUBJECT,notEmpty"`
	}
)

func New() (*Config, error) {

	loadDotEnv()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Manual validation if needed
	if cfg.Server.GRPCServer.Port == 0 {
		return nil, fmt.Errorf("GRPC_PORT must be set and non-zero")
	}

	return &cfg, nil
}

func loadDotEnv() error {
	filePath := fmt.Sprintf(".env")

	if _, err := os.Stat(filePath); err == nil {
		return godotenv.Load(filePath)
	}

	filePath = filepath.Join("..", fmt.Sprintf(".env"))
	return godotenv.Load(filePath)
}
