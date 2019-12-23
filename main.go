package main

import (
	"fmt"
	"log"
	"github.com/redisdemo/redisgo"
	"github.com/redisdemo/goredis"
	
)

func main(){
	fmt.Println("Go is working")
	log.Println("Log is working")	
	redisgo.Demo()
	goredis.Demo()

}