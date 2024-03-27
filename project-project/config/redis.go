package config

import (
	"github.com/go-redis/redis/v8"
	"test.com/project-project/internal/dao"
)

func (c *Config) ReConnRedis() {
	rdb := redis.NewClient(c.ReadRedisConfig())
	rc := &dao.RedisCache{
		Rdb: rdb,
	}
	dao.Rc = rc
}
