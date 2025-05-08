package database

import (
	config "TaskManagmentApis/configs"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService defines methods for accessing Redis
type RedisService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)                            // Add this if needed for Incr operation
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) // Add Expire method
	GetClient() *redis.Client
	Ping() error
	Close() error
}

// redisClient implements RedisService
type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService creates a new RedisService instance
func NewRedisService(ctx context.Context) (RedisService, error) {
	cfg := config.Config

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword, // Now from config
		DB:       0,
	})

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctxWithTimeout).Result()
	if err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	log.Println("✅ Successfully connected to Redis.")

	return &redisClient{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with an expiration time
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete store value from Redis
func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Incr increments a Redis value by 1
func (r *redisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Expire sets the expiration for a Redis key
func (r *redisClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, expiration).Result()
}

// GetClient returns the Redis client
func (r *redisClient) GetClient() *redis.Client {
	return r.client
}

// Ping checks Redis connection
func (r *redisClient) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// Close gracefully closes Redis connection
func (r *redisClient) Close() error {
	return r.client.Close()
}
