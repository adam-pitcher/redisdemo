package goredis

import (
    "log"
  	"github.com/go-redis/redis"
)

//Demo for using Go-Redis packages with Redis
func Demo(){
  log.Println("Starting basic Redis connection and transactions using go-redis library")

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	pong, err := client.Ping().Result()
	log.Println(pong,err)

	val, err := client.Get(")
}