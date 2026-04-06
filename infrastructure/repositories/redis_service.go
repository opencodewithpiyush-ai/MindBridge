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

func (r *RedisClient) CreateSession(token string, userID string, expiry time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Set(ctx, "session:"+token, userID, expiry).Err()
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

func (r *RedisClient) DeleteSession(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Del(ctx, "session:"+token).Err()
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

func (r *RedisClient) Close() error {
	return r.client.Close()
}
