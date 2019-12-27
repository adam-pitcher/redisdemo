package redisrest

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"strings"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"github.com/go-redis/redis/v7"
	"github.com/go-redis/redis_rate/v8"
	"github.com/gorilla/mux"
)

type demoData struct {
	ID	string `json:"ID"`
	Name string `json:"name"`
	Exp string `json:"exp"`
}

var limiter *redis_rate.Limiter

//Verifies a the redis demo is running and accepting requests
func verifyConn(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Redis Demo is available")
}

//Adds data to the demo redis database
func addDemoData(w http.ResponseWriter, r *http.Request){
	var data demoData	
	var apikey = r.Header.Get("apikey")
	var goClient *redis.Client	
	log.Printf("API KEY = %s ",apikey)
	if apikey == "TestKey1" {
	//Get a client from the pool and establish a limiter
	log.Println("Setting client with limiter 10 in a second")
		goClient = getGoClient(30, time.Second)
	} else {
		log.Println("Setting client with limiter 20 in a second")
		goClient = getGoClient(10, time.Second)
	}

	//Get the limiter status for the provided key
	limiterResult, err := limiter.Allow(apikey)
	if err != nil {
		panic(err)
	}	
	
	log.Printf("Limit and remaining: %d",limiterResult.Remaining)
	//If the request is within the limiter conditions, process the request
	if limiterResult.Allowed {

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w,"Error reading request data or request data invalid")
		}

		json.Unmarshal(reqBody,&data)
		
		data.ID = getUUID()
		duration := getDuration(data.Exp)
	
		marshaledJSONData,err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}	

		log.Println("Setting Redis Data")
		err = goClient.Set(data.ID, marshaledJSONData, duration).Err()
		if err !=nil{
				log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)

		encodeResult(w,data)	
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		encodeResult(w,"Too many requests")
	}
}

//Gets all keys currently stored in the redis database
func getKeys(w http.ResponseWriter, r *http.Request){	

	goClient := getGoClient(10, time.Minute)
	
	//Get the limiter status for the provided key
	limiterResult, err := limiter.Allow(r.Header.Get("apikey"))
	if err != nil {
		panic(err)
	}	
	
	if limiterResult.Allowed{
		result, err := goClient.Keys("*").Result()
		if err != nil{
			w.WriteHeader(http.StatusNotFound)
			encodeResult(w,"There was a problem retrieving keys or no keys exist")
			log.Panic(err)	
		}
		w.WriteHeader(http.StatusOK)
		encodeResult(w,result)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		encodeResult(w,"Too many requests")		
	}
}

//Gets data from the demo redis databased based on the provided key value
func getDemoData(w http.ResponseWriter, r *http.Request){
	goClient := getGoClient(10, time.Minute)

	//Get the limiter status for the provided key
	limiterResult, err := limiter.Allow(r.Header.Get("apikey"))
	if err != nil {
		panic(err)
	}	

	if limiterResult.Allowed{
		demoKey := mux.Vars(r)["id"]

		result,err := goClient.Get(demoKey).Result()
		if err != nil{
			w.WriteHeader(http.StatusNoContent)
			log.Printf("There was a problem retrieving data for id %s : %s",demoKey,err)
		} 
		w.WriteHeader(http.StatusOK)
		encodeResult(w,result)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		encodeResult(w,"Too many requests")
	}	
}

func deleteDemoData(w http.ResponseWriter, r *http.Request){
	goClient := getGoClient(10, time.Minute)

	limiterResult, err := limiter.Allow(r.Header.Get("apikey"))
	if err != nil {
		panic(err)
	}

	if limiterResult.Allowed{
		demoKey := mux.Vars(r)["id"]

		err := goClient.Del(demoKey).Err()

		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("There was a problem deleting key %s : %s",demoKey,err)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		encodeResult(w,"Too many requests")
	}	
}

func flushCurrentDB(w http.ResponseWriter, r *http.Request){
	goClient := getGoClient(10, time.Minute)

	limiterResult, err := limiter.Allow(r.Header.Get("apikey"))
	if err != nil {
		panic(err)
	}

	if limiterResult.Allowed{
		err := goClient.FlushDB().Err()

		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("There was a problem flushing the current db", err)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusOK)
		encodeResult(w,"Too many requests")
	}
}

//sets up a go client and limiter
func getGoClient(rate int, period time.Duration) *redis.Client {

	log.Println("Getting new redis client")

	poolTimeout,_ := time.ParseDuration("10m")

	goClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		PoolSize: 5,
		PoolTimeout: poolTimeout,
	})	

	limiter = redis_rate.NewLimiter(goClient,&redis_rate.Limit{
		Burst: rate,
		Rate: rate,
		Period: period,
	})	

	return goClient
} 

//Gets a GUID/UUID
func getUUID() string{
	out, _ := exec.Command("uuidgen").Output()
	return strings.TrimSuffix(string(out),"\n")
}

//Conerts a string to a duration
func getDuration(timeoutString string) time.Duration{
	duration,err := time.ParseDuration(timeoutString)
		if err != nil {
			log.Fatal(err)
		}
	return duration
}

//Json encodes the provided data to the response writer
func encodeResult(w http.ResponseWriter, data interface{}){
	json.NewEncoder(w).Encode(data)
}