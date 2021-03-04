package ginsessions

import (
	"gin-example/pkg/cache"
	"gin-example/pkg/sessions"
)

type redisStore struct {
	*sessions.RedisStore
}

// GinStoreInterface
type GinStoreInterface interface {
	// sessions.Store
	sessions.Store

	// set http.cookie options default max age
	SetMaxAge(maxAge int)
	// set http.cookie options max length
	SetMaxLength(maxLength int)
	// set store redis key prefix
	SetRedisKeyPrefix(prefix string)
	// set SetSerializer method
	SetSerializer(i sessions.SessionSerializerInterface)
}

func NewRedisStore(keyPairs ...[]byte) (GinStoreInterface, error) {
	rs, err := sessions.NewRedisStore(cache.GetSessionCache(), keyPairs...)
	if err != nil {
		return nil, err
	}
	return &redisStore{RedisStore: rs}, nil
}

func (c *redisStore) CovertOptions(options Options) {
	c.RedisStore.Options = options.ToOptions()
}
