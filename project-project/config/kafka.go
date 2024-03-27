package config

import (
	"context"
	"go.uber.org/zap"
	"test.com/project-common/kk"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/repo"
	"time"
)

var kw *kk.KafkaWriter

func InitKafkaWriter() func() {
	kw = kk.GetWriter("localhost:9092")
	return kw.Close
}

func SendLog(data []byte) {
	kw.Send(kk.LogData{
		Topic: "msproject_log",
		Data:  data,
	})
}

func SendCache(data []byte) {
	kw.Send(kk.LogData{
		Topic: "msproject_cache",
		Data:  data,
	})
}

type KafkaCache struct {
	R     *kk.KafkaReader
	cache repo.Cache
}

func (c *KafkaCache) DeleteCache() {
	for {
		message, err := c.R.R.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("DeleteCache ReadMessage err", zap.Error(err))
			continue
		}
		zap.L().Info("收到缓存", zap.String("value", string(message.Value)))
		if "task" == string(message.Value) {
			fields, err := c.cache.HKeys(context.Background(), "task")
			if err != nil {
				zap.L().Error("DeleteCache HKeys err", zap.Error(err))
				continue
			}
			time.Sleep(1 * time.Second)
			c.cache.Delete(context.Background(), fields)
		}
	}
}

func NewCacheReader() *KafkaCache {
	reader := kk.GetReader([]string{"localhost:9092"}, "cache_group", "msproject_cache")
	return &KafkaCache{
		R:     reader,
		cache: dao.Rc,
	}
}
