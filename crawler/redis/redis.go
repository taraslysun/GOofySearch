package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	redisClient *redis.Client
)

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func CheckTime(link string) (time.Duration, error) {
	InitRedis()
	ttl := redisClient.TTL(link)
	return ttl.Val(), ttl.Err()
}

func SetVal(link string, val string) {
	InitRedis()
	err := redisClient.Set(link, val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func SetTime(link string, duration time.Duration) {
	InitRedis()
	err := redisClient.Expire(link, time.Second*duration).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func IsLinkVisited(link string) bool {
	InitRedis()
	val, err := redisClient.Get(link).Result()
	if err != nil {
		fmt.Println(err)
	}
	if val == "true" {
		return true
	}
	return false
}
