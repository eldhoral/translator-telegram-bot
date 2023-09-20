package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error in load env file: %v\n ", err)
	}
}

func initRedisClient() (client *redis.Client, err error) {
	var ctx = context.TODO()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username: os.Getenv("REDIS_USERNAME"),
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
	redisClient, _ := initRedisClient()
	go TelegramBot(redisClient)
	// Init router
	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheck).Methods("GET")
	// start app
	startApp(r)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up :)")
}

func startApp(r *mux.Router) {
	http.ListenAndServe(":8085", r)
}
