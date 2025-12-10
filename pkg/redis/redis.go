package redis

import (
	"context"
	"fmt"
	"time"

	"suitemedia/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client    *redis.Client
	keyPrefix string
}

func NewClient(cfg config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		client:    client,
		keyPrefix: cfg.KeyPrefix,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, c.keyPrefix+key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, c.keyPrefix+key, value, expiration).Err()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.keyPrefix+key).Err()
}

func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
