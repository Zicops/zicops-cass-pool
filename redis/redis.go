package redis

import (
	"fmt"
	"os"
	"time"

	"context"

	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	redisClient *redis.Client
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Initialize() (*redis.Client, error) {
	redisHost := getEnv("REDIS_HOST", "127.0.0.1")
	logrus.Info("Redis Host: ", redisHost)
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})
	redisClient = rdb
	return rdb, nil
}

func CloseRedisConn(conn redis.Conn) {
	conn.Close()
}

func GetRedisValue(ctx context.Context, key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

func SetRedisValue(ctx context.Context, key string, value string) error {
	return redisClient.Set(ctx, key, value, 0).Err()
}

func SetTTL(ctx context.Context, key string, ttl int) error {
	return redisClient.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
}
