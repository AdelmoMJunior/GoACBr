package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/AdelmoMJunior/GoACBr/internal/config"
)

// Cache interface defines the common cache operations.
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
	IsAvailable() bool
}

// RedisCache implementation of Cache using go-redis.
type RedisCache struct {
	client *redis.Client
	isDown bool
}

// NewRedisCache initializes a new Redis cache connection.
func NewRedisCache(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("Failed to connect to Redis. Running in degraded mode without cache.", "error", err)
		return &RedisCache{client: client, isDown: true}, nil
	}

	slog.Info("Connected to Redis successfully", "host", cfg.Host)
	return &RedisCache{client: client, isDown: false}, nil
}

// IsAvailable checks if the Redis connection is active.
func (c *RedisCache) IsAvailable() bool {
	return !c.isDown
}

// Get retrieves and unmarshals a value from cache.
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	if c.isDown {
		return redis.Nil // simulate cache miss
	}

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return err
		}
		slog.Error("Redis GET error", "key", key, "error", err)
		c.checkConnection(ctx)
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set marshals and stores a value in cache with a TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if c.isDown {
		return nil
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	err = c.client.Set(ctx, key, bytes, ttl).Err()
	if err != nil {
		slog.Error("Redis SET error", "key", key, "error", err)
		c.checkConnection(ctx)
		return err
	}

	return nil
}

// Delete removes a key from cache.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if c.isDown {
		return nil
	}

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		slog.Error("Redis DEL error", "key", key, "error", err)
		c.checkConnection(ctx)
	}
	return err
}

// Close closes the Redis connection.
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// checkConnection performs a soft check to see if Redis came back or went down.
func (c *RedisCache) checkConnection(ctx context.Context) {
	err := c.client.Ping(ctx).Err()
	if err != nil && !c.isDown {
		slog.Warn("Redis connection lost. Switching to degraded mode.")
		c.isDown = true
	} else if err == nil && c.isDown {
		slog.Info("Redis connection recovered. Enabling cache.")
		c.isDown = false
	}
}
