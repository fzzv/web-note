package gredis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/fzzv/go-gin-example/pkg/setting"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// Setup 初始化 Redis 客户端
func Setup() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         setting.RedisSetting.Host,     // e.g. "localhost:6379"
		Password:     setting.RedisSetting.Password, // "" if no password
		DB:           0,                             // 默认数据库
		PoolSize:     setting.RedisSetting.MaxActive,
		MinIdleConns: setting.RedisSetting.MaxIdle,
	})

	// 测试连接
	result, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("PING 返回:", result)
	return nil
}

// Set 设置 key 并指定过期时间（秒）
func Set(key string, data interface{}, expireSeconds int) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, time.Duration(expireSeconds)*time.Second).Err()
}

// Exists 判断 key 是否存在
func Exists(key string) bool {
	ok, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return ok > 0
}

// Get 获取 key
func Get(key string) ([]byte, error) {
	val, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // key 不存在
	}
	return val, err
}

// Delete 删除 key
func Delete(key string) (bool, error) {
	deleted, err := rdb.Del(ctx, key).Result()
	return deleted > 0, err
}

// LikeDeletes 按模式删除（模糊匹配）
func LikeDeletes(pattern string) error {
	iter := rdb.Scan(ctx, 0, "*"+pattern+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		_, err := rdb.Del(ctx, key).Result()
		if err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
