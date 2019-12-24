package redisrest

import (
	"fmt"
	"strings"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
)

type redisData struct {
	ID	string `json:"ID"`
	Name string `json:"Name"`
}

func InitializeServer(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/",verifyConn)
	router.HandleFunc("/addRedisData",addRedisData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080",router))
}

func verifyConn(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Redis Demo is available")
}

func addRedisData(w http.ResponseWriter, r *http.Request){
	var data redisData
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w,"Error reading request data or request data invalid")
	}

	json.Unmarshal(reqBody,&data)

	out, err := exec.Command("uuidgen").Output()
	data.ID = strings.TrimSuffix(string(out),"\n")

	log.Println("Getting new redis client")

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	jsn,err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	err = client.Set("1", jsn, 0).Err()
	if err !=nil{
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
	
	
}
