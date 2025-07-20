package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	dbIndex := getRedisDB()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbIndex,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully!")
}

func getRedisDB() int {
	dbIndex, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return 0
	}
	return dbIndex
}
