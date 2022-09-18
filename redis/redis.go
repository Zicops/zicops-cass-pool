package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	//"github.com/garyburd/redigo/redis"
)

var (
	redisPool *redis.Pool
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Initialize() (*redis.Pool, error) {
	redisHost := getEnv("REDIS_HOST", "127.0.0.1")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisPool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
			if err != nil {
				return nil, err
			}
			if redisPassword != "" {
				if _, err := c.Do("AUTH", redisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
	return redisPool, nil
}

func GetRedisPool() *redis.Pool {
	return redisPool
}

func GetRedisConn() redis.Conn {
	return redisPool.Get()
}

func CloseRedisConn(conn redis.Conn) {
	conn.Close()
}

func GetRedisValue(key string) (string, error) {
	conn := GetRedisConn()
	defer CloseRedisConn(conn)
	return redis.String(conn.Do("GET", key))
}

func SetRedisValue(key string, value string) error {
	conn := GetRedisConn()
	defer CloseRedisConn(conn)
	_, err := conn.Do("SET", key, value)
	return err
}

func DeleteRedisValue(key string) error {
	conn := GetRedisConn()
	defer CloseRedisConn(conn)
	_, err := conn.Do("DEL", key)
	return err
}
