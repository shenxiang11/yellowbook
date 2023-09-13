package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"yellowbook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Set(ctx context.Context, u domain.User) error
	Get(ctx context.Context, id uint64) (domain.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Get(ctx context.Context, id uint64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *RedisUserCache) Delete(ctx context.Context, id uint64) error {
	key := cache.key(id)
	return cache.client.Del(ctx, key).Err()
}

func (cache *RedisUserCache) key(id uint64) string {
	return fmt.Sprintf("user:info:%d", id)
}
