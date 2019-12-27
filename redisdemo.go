package main

import (
	"fmt"
	"log"
	"flag"
	"github.com/redisdemo/redisgo"
	"github.com/redisdemo/goredis"
	"github.com/redisdemo/redisrest"	
)

func main(){
	redisgoFlag := flag.Bool("redisgo",false,"Run redis-go demo")
	goredisFlag := flag.Bool("goredis",false,"Run go-redis demo")
	goredisRestFlag := flag.Bool("goredissrv",false,"Start a server which interacts with Redis via a REST Api using the goredis packages")
	redisgoRestFlag := flag.Bool("redisgorsrv",false,"Start a server which interacts with Redis via a REST Api using the redisgo packages")

	flag.Parse()	

	if *redisgoFlag{
		redisgo.Demo()
	} else if *goredisFlag{
		goredis.Demo()
	} else if *goredisRestFlag{
		log.Println("Initializing Server")
		redisrest.InitializeServer()
	} else if *redisgoRestFlag{
		redisrest.InitializeServer()
	} else {
		fmt.Println("Please select a run option for the demo:")
		fmt.Println("-redisgo to run an automated demo using redigo packages")
		fmt.Println("-goredis to run an automated demo using goredis packages")
		fmt.Println("-goredissrv to launch a REST API server which processes requests to Redis")
	}
}
