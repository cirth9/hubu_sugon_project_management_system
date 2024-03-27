package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rc *RedisCache

type RedisCache struct {
	Rdb *redis.Client
}

//func init() {
//	Rdb := redis.NewClient(config.C.ReadRedisConfig())
//	Rc = &RedisCache{
//		Rdb: Rdb,
//	}
//}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.Rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.Rdb.Get(ctx, key).Result()
	return result, err
}

func (rc *RedisCache) HSet(ctx context.Context, key string, field string, value string) {
	rc.Rdb.HSet(ctx, key, field, value)
}

func (rc *RedisCache) HKeys(ctx context.Context, key string) ([]string, error) {
	result, err := rc.Rdb.HKeys(ctx, key).Result()
	return result, err
}

func (rc *RedisCache) Delete(ctx context.Context, keys []string) {
	rc.Rdb.Del(ctx, keys...)
}
