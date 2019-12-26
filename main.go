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

	fmt.Println("Go is working")
	log.Println("Log is working")	

	if *redisgoFlag{
		redisgo.Demo()
	}

	if *goredisFlag{
		goredis.Demo()
	}

	if *goredisRestFlag{
		redisrest.InitializeServer()
	}

	if *redisgoRestFlag{
		redisrest.InitializeServer()
	}
}

