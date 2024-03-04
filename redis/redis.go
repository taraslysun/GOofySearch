package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	redisClient *redis.Client
)

func redis_init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func check_time(link string) (time.Duration, error) {
	redis_init()
	ttl := redisClient.TTL(link)
	return ttl.Val(), ttl.Err()
}

func set_time(link string, duration time.Duration) {
	redis_init()
	err := redisClient.Expire(link, time.Second*duration).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func is_link_visited(link string) bool {
	redis_init()
	val, err := redisClient.Get(link).Result()
	if err != nil {
		fmt.Println(err)
	}
	if val == "true" {
		return true
	}
	return false
}
