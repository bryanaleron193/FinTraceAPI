package cache

import (
	"context"
	"log"
	"simple-gin-backend/internal/config"
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
		Addr:     config.AppConfig.RedisAddr,
		Password: config.AppConfig.RedisPassword,
		DB:       dbIndex,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully!")
}

func getRedisDB() int {
	dbIndex, err := strconv.Atoi(config.AppConfig.RedisDB)
	if err != nil {
		return 0
	}
	return dbIndex
}
