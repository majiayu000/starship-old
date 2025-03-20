package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/majiayu000/cc-starship/pkg/config"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// RedisDB represents a Redis database connection
type RedisDB struct {
	*redis.Client
	logger *logger.Logger
}

// NewRedisDB creates a new Redis database connection
func NewRedisDB(cfg *config.Config, logger *logger.Logger) (*RedisDB, error) {
	// Check if Redis is enabled in the configuration
	if !cfg.Databases.Redis.Enabled {
		return nil, fmt.Errorf("redis is not enabled in configuration")
	}

	// Create Redis client options
	opts := &redis.Options{
		Addr:         cfg.Databases.Redis.Address,
		Password:     cfg.Databases.Redis.Password,
		DB:           cfg.Databases.Redis.DB,
		PoolSize:     cfg.Databases.Redis.PoolSize,
		MinIdleConns: cfg.Databases.Redis.MinIdleConns,
		MaxRetries:   cfg.Databases.Redis.MaxRetries,
		DialTimeout:  cfg.Databases.Redis.DialTimeout,
	}

	// Create Redis client
	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis database")

	return &RedisDB{
		Client: client,
		logger: logger,
	}, nil
}

// Close closes the Redis connection
func (r *RedisDB) Close() error {
	r.logger.Info("Closing Redis database connection")
	return r.Client.Close()
}

// HealthCheck checks the Redis connection health
func (r *RedisDB) HealthCheck(ctx context.Context) error {
	return r.Ping(ctx).Err()
}

// Type returns the database type
func (r *RedisDB) Type() string {
	return "redis"
}

// Set sets a key-value pair in Redis with an optional expiration
func (r *RedisDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value by key from Redis
func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete deletes a key from Redis
func (r *RedisDB) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisDB) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	return result > 0, err
}

// Increment increments a key's value
func (r *RedisDB) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// HSet sets hash fields in Redis
func (r *RedisDB) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	return r.Client.HSet(ctx, key, values).Err()
}

// HGet gets a hash field from Redis
func (r *RedisDB) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll gets all hash fields from Redis
func (r *RedisDB) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// Expire sets an expiration on a key
func (r *RedisDB) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// Publish publishes a message to a Redis channel
func (r *RedisDB) Publish(ctx context.Context, channel, message string) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to Redis channels
func (r *RedisDB) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channels...)
}
