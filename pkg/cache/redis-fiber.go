package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	RedisClient *redis.Client
	Ctx         context.Context
}

func (rs RedisStore) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val, err := rs.RedisClient.Get(rs.Ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, nil
}

func (rs RedisStore) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return rs.RedisClient.Set(context.Background(), key, val, exp).Err()
}

func (rs RedisStore) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	return rs.RedisClient.Del(context.Background(), key).Err()
}

func (rs RedisStore) Reset() error {
	return rs.RedisClient.FlushDB(context.Background()).Err()
}

func (rs RedisStore) Close() error {
	return rs.RedisClient.Close()
}
