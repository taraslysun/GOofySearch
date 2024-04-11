package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func newRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

func newCtx() context.Context {
	ctx := context.Background()
	return ctx
}

func addLink(client *redis.Client, ctx context.Context, link string, queue string) {
	err := client.LPush(ctx, queue, link).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func popLink(client *redis.Client, ctx context.Context, queue string) string {
	link, err := client.LPop(ctx, queue).Result()
	if err != nil {
		log.Fatal(err)
	}
	return link
}
