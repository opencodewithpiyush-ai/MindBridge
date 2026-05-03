package repositories

import (
	"context"
	"fmt"
	"mindbridge/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Username: config.RedisUsername,
		Password: config.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("[Redis] Connected to Redis")
	return &RedisClient{client: client}, nil
}

func (r *RedisClient) CreateSession(jti string, userID string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store session keyed by jti
	if err := r.client.Set(ctx, "session:"+jti, userID, expiry).Err(); err != nil {
		return err
	}
	// Also add jti to the user's session set
	return r.client.SAdd(ctx, "user_sessions:"+userID, jti).Err()
}

func (r *RedisClient) GetSession(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.client.Get(ctx, "session:"+token).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisClient) DeleteSession(jti string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve userID before deleting session
	userID, err := r.client.Get(ctx, "session:"+jti).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// Delete session key
	if err := r.client.Del(ctx, "session:"+jti).Err(); err != nil {
		return err
	}

	// Remove jti from user's session set
	if userID != "" {
		if err := r.client.SRem(ctx, "user_sessions:"+userID, jti).Err(); err != nil {
			// Non-critical, log but don't fail
			fmt.Printf("[Redis] Warning: could not remove jti from user set: %v\n", err)
		}
	}
	return nil
}

func (r *RedisClient) IsSessionValid(token string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.client.Exists(ctx, "session:"+token).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisClient) Incr(key string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisClient) Expire(key string, expiry time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r.client.Expire(ctx, key, expiry)
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
