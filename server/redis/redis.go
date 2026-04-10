package redisServer

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var client *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}

	log.Println("Redis 连接成功")
}

// GetRedis 获取 Redis 客户端
func GetRedis() *redis.Client {
	return client
}
