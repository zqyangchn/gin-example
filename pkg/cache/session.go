package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gin-example/pkg/logging"
	"gin-example/pkg/setting"
)

type SessionCacheRedisClientInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

var sessionCache SessionCacheRedisClientInterface

func Setup() error {
	// 初始化 sessionCache
	switch setting.SessionRedisSetting.Type {
	case "singlePoint":
		sessionCache = redis.NewClient(&redis.Options{
			Addr:     setting.SessionRedisSetting.Address,
			Password: setting.SessionRedisSetting.Password,
		})
	case "cluster":
		sessionCache = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          setting.SessionRedisSetting.Addresses,
			Password:       setting.SessionRedisSetting.Password,
			RouteRandomly:  false,
			RouteByLatency: false,
		})
	default:
		return errors.Errorf("unknown SessionRedisSetting.type: %s, use singlePoint|cluster", setting.SessionRedisSetting.Type)
	}
	if result, err := sessionCache.Ping(context.Background()).Result(); err != nil || result != "PONG" {
		return err
	}
	logging.Logger.Info("initialization sessionCache ok.", zap.String("mode", setting.SessionRedisSetting.Type))

	return nil
}

func GetSessionCache() SessionCacheRedisClientInterface {
	return sessionCache
}
