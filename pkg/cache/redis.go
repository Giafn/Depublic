package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Giafn/Depublic/configs"
	"github.com/redis/go-redis/v9"
)

func InitRedis(config *configs.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("redis connection failed")
	}

	return client, nil
}

type Cacheable interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) string
}

type cacheable struct {
	client *redis.Client
}

func NewCacheable(client *redis.Client) *cacheable {
	return &cacheable{
		client: client,
	}
}

func (c *cacheable) Set(key string, value interface{}, expiration time.Duration) error {
	operation := c.client.Set(context.Background(), key, value, expiration)
	if err := operation.Err(); err != nil {
		return err
	}
	return nil
}

func (c *cacheable) Get(key string) string {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return val
}
