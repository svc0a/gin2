package idempotency

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration // 设置键的过期时间
}

// NewRedisStore 构造函数，初始化 RedisStore
func NewRedisStore(addr string, password string, db int, ttl time.Duration) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 没有密码时使用 ""
		DB:       db,       // Redis 数据库编号
	})

	return &RedisStore{
		client: client,
		ttl:    ttl,
	}
}

// Store 方法实现，向 Redis 中存储键值对，并设置过期时间
func (r *RedisStore) Store(key string, value []byte) {
	data := string(value)
	// 将值序列化为字符串或字节数据
	err := r.client.Set(context.Background(), key, data, r.ttl).Err()
	if err != nil {
		// 这里可以根据需要处理错误，比如日志记录
	}
}

// Load 方法实现，从 Redis 中加载键值对
func (r *RedisStore) Load(key string) ([]byte, error) {
	value, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		// 如果键不存在，返回自定义错误
		return nil, errors.New("record not found")
	} else if err != nil {
		// 其他 Redis 错误
		return nil, err
	}
	return []byte(value), nil
}
