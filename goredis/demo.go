package goredis

import (
    "log"
	  "github.com/go-redis/redis"
	  "time"
)

type demoClient struct {
	client *redis.Client
}

//Demo for using Go-Redis packages with Redis
func Demo(){

	log.Println("Starting basic Redis interaction using go-redis library")

	dClient := getRedisClient()  

	dClient.verifyRedisClient()
  
	dClient.setRedisData("name","GoRedisPackage",0)

	value := dClient.getRedisData("Name")

	log.Printf("Value is %s of type %T",value, value)

}

func getRedisClient() *demoClient {
	log.Println("Getting new redis client")

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	dClient := new(demoClient)
	dClient.client = client
	return dClient
}

func (dc demoClient) verifyRedisClient(){
	log.Println("Verifying connectivity with Redis")  

	pong, err := dc.client.Ping().Result()
	log.Println(pong,err)
}

func (dc demoClient) setRedisData(key string, value string, expiration time.Duration) {
//Set a value in Redis  Set(Key, Value, Expiration(0=none))
	err := dc.client.Set(key, value, expiration).Err()
	if err !=nil{
		log.Fatal(err)
	}
}

func (dc demoClient) getRedisData(key string) string{
	value, err := dc.client.Get("name").Result()
	if err != nil{
		log.Fatal(err)
	}
	return value
}