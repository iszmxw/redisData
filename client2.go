package main

import (
	"fmt"
	"log"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/huobi"
	"redisData/setting"
)

func main() {
	symbol := "btcusdt"
	period := "1min"
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
	//go QueryKlineData()
	fmt.Println("success")
	huobi.NewSubscribe(symbol, period)
}
