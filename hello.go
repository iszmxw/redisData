package main

import (
	"fmt"
	"log"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/setting"
)

func main() {
	//初始化viper
	if err := setting.Init("");err !=nil{
		log.Println("viper init fail")
		return
	}

	//初始化redis
	if err := redis.InitClient();err != nil{
		fmt.Printf("init redis fail err:%v/n", err)
		return
	}

	var symbols = []string{"btcusdt"}

	if err := logic.AutoGetRedisData(symbols);err !=nil{
		return
	}
}