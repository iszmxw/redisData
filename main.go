package main

import (
	"fmt"
	"log"
	"redisData/dao/redis"
	"redisData/routes"
	"redisData/setting"
)

//var symbol = make([]string, 0)

//  QueryKlineData 每隔30s发送一次请求，请求全部类型的交易对
//  如果没有key就创建，存在key就更新

func QueryKlineData() {
	//response, err := http.Get("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt")

}

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

	//初始化routes
	r := routes.SetUp()
	r.Run(":8887")
}
