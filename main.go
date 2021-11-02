package main

import (
	"fmt"
	"log"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/routes"
	"redisData/setting"
	"time"
)

//var symbol = make([]string, 0)

//  QueryKlineData 每隔30s发送一次请求，请求全部类型的交易对
//  如果没有key就创建，存在key就更新

func QueryKlineData() {
	//response, err := http.Get("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt")
	if err := logic.AutoGetRedisData(); err != nil {
		fmt.Printf("AutoGetRedisData is fail %v", err)
		return
	}
	time.AfterFunc(30*time.Second, QueryKlineData)
}

func main() {
	//初始化viper
	if err := setting.Init(""); err != nil {
		log.Println("viper init fail")
		return
	}

	//初始化MySQL
	if err := mysql.InitMysql(); err != nil {
		fmt.Printf("init mysql fail err:%v/n", err)
		return
	}
	defer mysql.Close()

	//初始化redis
	if err := redis.InitClient(); err != nil {
		fmt.Printf("init redis fail err:%v/n", err)
		return
	}
	defer redis.Close()
	//起协程每隔30s请求一次并且缓存到redis
	//var symbols = []string{"btcusdt"}
	//go logic.AutoGetRedisData()
	//if err := logic.AutoGetRedisData(); err != nil {
	//	fmt.Printf("AutoGetRedisData is fail %v", err)
	//	return
	//}
	go QueryKlineData()
	fmt.Println("success")
	//初始化routes
	r := routes.SetUp()
	r.Run(":8887")
}
