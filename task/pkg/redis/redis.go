package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Config struct {
	Host         string
	Password     string
	TLSEnable    bool
	DialTimeout  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func GetRedisConfig() Config {
	dialTimeout, _ := time.ParseDuration(os.Getenv("REDIS_DIAL_TIMEOUT"))
	writeTimeout, _ := time.ParseDuration(os.Getenv("REDIS_WRITE_TIMEOUT"))
	readTimeout, _ := time.ParseDuration(os.Getenv("REDIS_READ_TIMEOUT"))

	tlsEnable := os.Getenv("REDIS_TLS_ENABLE") == "true"

	return Config{
		Host:         os.Getenv("REDIS_HOSTS"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		TLSEnable:    tlsEnable,
		DialTimeout:  dialTimeout,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

type Client struct {
	client *redis.Client
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	options := &redis.Options{
		Addr:         cfg.Host,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Password:     cfg.Password,
	}

	if cfg.TLSEnable {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(options)

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{client}, nil
}

func (c *Client) Close() error {
	err := c.client.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}

func (c *Client) Ping(ctx context.Context) error {
	err := c.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis: %w", err)
	}

	return nil
}

func (c *Client) Unwrap() *redis.Client {
	return c.client
}
