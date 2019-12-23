package redisgo

import (
	"fmt"
	"log"
	"github.com/gomodule/redigo/redis"
)

// Demo code for using RedisGo packages with Redis
func Demo(){
	fmt.Println("GO is working")

	log.Println("Starting basic Redis connection and transactions using RedisGo library")
	conn, err := redis.Dial("tcp","localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Dial success ")

	_,err = conn.Do("HMSET","Name:1", "First","AdamP")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("AdamP Added")

	first,err := redis.String(conn.Do("HGET","Name:1","First"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Name = %s",first)
}