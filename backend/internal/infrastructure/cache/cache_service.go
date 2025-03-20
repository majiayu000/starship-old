package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/majiayu000/cc-starship/internal/infrastructure/database"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// CacheService provides caching functionality
type CacheService struct {
	redis       *database.RedisDB
	logger      *logger.Logger
	memoryCache map[string]cacheItem
	mu          sync.RWMutex
	useRedis    bool
}

// cacheItem represents an item in the memory cache
type cacheItem struct {
	Value      []byte
	Expiration time.Time
}

// NewCacheService creates a new cache service
func NewCacheService(dbManager *database.DataAccessManager, logger *logger.Logger) *CacheService {
	service := &CacheService{
		logger:      logger,
		memoryCache: make(map[string]cacheItem),
		useRedis:    false,
	}

	// Use Redis if available
	if dbManager.Redis != nil {
		service.redis = dbManager.Redis
		service.useRedis = true
		logger.Info("Cache service initialized with Redis backend")
	} else {
		logger.Info("Cache service initialized with in-memory backend")
	}

	return service
}

// Set stores a value in the cache
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Marshal value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	if c.useRedis {
		// Use Redis
		return c.redis.Set(ctx, key, jsonValue, expiration)
	}

	// Use memory cache
	c.mu.Lock()
	defer c.mu.Unlock()

	c.memoryCache[key] = cacheItem{
		Value:      jsonValue,
		Expiration: time.Now().Add(expiration),
	}
	return nil
}

// Get retrieves a value from the cache
func (c *CacheService) Get(ctx context.Context, key string, result interface{}) (bool, error) {
	var jsonValue []byte
	var found bool

	if c.useRedis {
		// Use Redis
		value, err := c.redis.Get(ctx, key)
		if err != nil {
			// Check if key not found
			if err.Error() == "redis: nil" {
				return false, nil
			}
			return false, err
		}
		jsonValue = []byte(value)
		found = true
	} else {
		// Use memory cache
		c.mu.RLock()
		item, exists := c.memoryCache[key]
		c.mu.RUnlock()

		// Check expiration
		if !exists || time.Now().After(item.Expiration) {
			if exists {
				// Clean up expired item
				c.mu.Lock()
				delete(c.memoryCache, key)
				c.mu.Unlock()
			}
			return false, nil
		}

		jsonValue = item.Value
		found = true
	}

	// Unmarshal JSON into result
	if found {
		if err := json.Unmarshal(jsonValue, result); err != nil {
			return true, fmt.Errorf("failed to unmarshal cache value: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// Delete removes a value from the cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	if c.useRedis {
		// Use Redis
		return c.redis.Delete(ctx, key)
	}

	// Use memory cache
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.memoryCache, key)
	return nil
}

// Clear removes all values from the cache
func (c *CacheService) Clear(ctx context.Context) error {
	if c.useRedis {
		c.logger.Warn("Clear operation not implemented for Redis cache")
		return nil
	}

	// Use memory cache
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memoryCache = make(map[string]cacheItem)
	return nil
}

// SetWithPrefix stores a value with a prefix
func (c *CacheService) SetWithPrefix(ctx context.Context, prefix, key string, value interface{}, expiration time.Duration) error {
	fullKey := fmt.Sprintf("%s:%s", prefix, key)
	return c.Set(ctx, fullKey, value, expiration)
}

// GetWithPrefix retrieves a value with a prefix
func (c *CacheService) GetWithPrefix(ctx context.Context, prefix, key string, result interface{}) (bool, error) {
	fullKey := fmt.Sprintf("%s:%s", prefix, key)
	return c.Get(ctx, fullKey, result)
}

// DeleteWithPrefix removes a value with a prefix
func (c *CacheService) DeleteWithPrefix(ctx context.Context, prefix, key string) error {
	fullKey := fmt.Sprintf("%s:%s", prefix, key)
	return c.Delete(ctx, fullKey)
}
