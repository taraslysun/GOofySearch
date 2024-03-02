package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	redisClient *redis.Client
)

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func isLinkAccessAllowed(link string, timeout time.Duration) bool {
	val, err := redisClient.Get(link).Result()
	if errors.Is(err, redis.Nil) {
		return true
	} else if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	lastAccessedTime, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	return time.Since(lastAccessedTime) > timeout
}

func markLinkAsAccessed(link string, timeout time.Duration) error {
	expiration := timeout // Set expiration time for the link
	return redisClient.Set(link, time.Now().Format(time.RFC3339Nano), expiration).Err()
}
