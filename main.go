package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panic(err)
	}
}

func initRedisClient() (client *redis.Client, err error) {
	var ctx = context.TODO()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_PASSWORD"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       cast.ToInt(os.Getenv("REDIS_DB")),
	})
	err = client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error in initializing redis client: " + err.Error())
		return
	}
	log.Println("Redis is running")
	return
}

func main() {
	loadConfig()
	redisClient, _ := initRedisClient()
	TelegramBot(redisClient)
}
