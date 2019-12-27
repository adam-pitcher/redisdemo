package redisrest

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

//InitializeServer - Sets up and intializes a server to handle incomming requests
func InitializeServer(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/",verifyConn)
	router.HandleFunc("/addDemoData",addDemoData).Methods("POST")
	router.HandleFunc("/getDemoKeys",getKeys).Methods("GET")
	router.HandleFunc("/getDemoData/{id}",getDemoData).Methods("GET")
	router.HandleFunc("/delDemoData/{id}",deleteDemoData).Methods("DELETE")
	router.HandleFunc("/flushDemoDb",flushCurrentDB).Methods("DELETE")
	log.Println("Redis demo is now live and ready to process requests")
	log.Fatal(http.ListenAndServe(":8080",router))
}

