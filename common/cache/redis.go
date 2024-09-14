package cache

import (
	"context"
	"fmt"
	"ginProject/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var RC *RedisClient

func Init() {
	RC = GetClient()
}

func GetClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         config.RedisHost,
		Password:     config.RedisPwd,
		MaxRetries:   config.RedisMaxRetries,
		PoolSize:     config.RedisPoolSize,
		MinIdleConns: config.RedisMinIdleConns,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("redis ping err:%v", err))
	}

	return &RedisClient{client: client, ctx: ctx}
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) HSet(key string, field string, value interface{}) error {
	return r.client.HSet(r.ctx, key, field, value).Err()
}

func (r *RedisClient) HGet(key string, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

func (r *RedisClient) HDel(key string, field string) error {
	return r.client.HDel(r.ctx, key, field).Err()
}

func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

func (r *RedisClient) HKeys(key string) ([]string, error) {
	return r.client.HKeys(r.ctx, key).Result()
}
