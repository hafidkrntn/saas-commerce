package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Cache interface {
	Del(ctx context.Context, keys ...string) error
	Fetch(ctx context.Context, key string, ttl time.Duration, dest any, loader func() (any, error)) error
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}

type RedisCache struct {
	client *redis.Client
	log    *logrus.Logger
}

func NewRedisCache(log *logrus.Logger) (Cache, func(), error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.WithError(err).Fatal("redis connection failed")
		return nil, nil, err
	}

	log.Info("redis connected")

	cleanup := func() {
		err := client.Close()
		if err != nil {
			log.WithError(err).Error("failed to close redis client")
			return
		}
		log.Info("redis client closed")
	}

	return &RedisCache{client: client, log: log}, cleanup, nil
}

// Del removes one or more keys.
func (c *RedisCache) Del(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Fetch is a cache-aside helper. It tries to get the value from cache;
// on a miss it calls loader(), stores the result, and returns it.
// If Redis is unavailable the loader is called directly without caching.
func (c *RedisCache) Fetch(
	ctx context.Context,
	key string,
	ttl time.Duration,
	dest any,
	loader func() (any, error),
) error {
	if err := c.Get(ctx, key, dest); err == nil {
		return nil
	} else if err != redis.Nil {
		value, err := loader()
		if err != nil {
			return err
		}
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(value))
		return nil
	}

	value, err := loader()
	if err != nil {
		return err
	}

	reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(value))

	if err := c.Set(ctx, key, value, ttl); err != nil {
		c.log.WithError(err).Errorf("cache.Fetch set error (key=%s)", key)
	}

	return nil
}

// Get deserializes the cached JSON into dest. Returns redis.Nil if not found.
func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	b, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err // caller checks redis.Nil
	}
	return json.Unmarshal(b, dest)
}

// Set serializes value as JSON and stores it with a TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache.Set marshal: %w", err)
	}
	return c.client.Set(ctx, key, b, ttl).Err()
}
