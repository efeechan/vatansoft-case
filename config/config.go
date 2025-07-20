package config

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var REDIS *redis.Client

func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitRedis() {
	REDIS = redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_HOST", "localhost:6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}
