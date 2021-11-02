package main

import (
	"fmt"
	"log"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/logic"
	"redisData/setting"
)

func main() {
	//初始化viper
	if err := setting.Init(""); err != nil {
		log.Println("viper init fail")
		return
	}
	//初始化msyql
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

	//var symbols = []string{"btcusdt", "etcusdt", "xrpusdt", "adausdt", "ltcusdt,", "xemusdt", "dashusdt", "xlmusdt", "ethbtc", "ethbtc", "eosbtc", "dotbtc", "dotbtc", "linketh", "adaeth", "xmreth"}

	if err := logic.AutoGetRedisData(); err != nil {
		fmt.Printf("AutoGetRedisData is fail %v", err)
		return
	}
	fmt.Println("success")
	//data, _ := logic.GetDataByKey("btcusdt")
	//fmt.Printf("btcusdt is %v", data)
	//go function()
	//func function(){
	//	// TODO 具体逻辑
	//
	//	// 每5分钟执行一次
	//	time.AfterFunc(5*time.Minute, function)
	//}

	//symbol, err := mysql.GetAllSymbol()
	//if err != nil {
	//	fmt.Printf("mysql.GetAllSymbol fail %v", err)
	//	return
	//}
	//ss := make([]string, 0, len(*symbol))
	//for _, value := range *symbol {
	//	ss = append(ss, value.Name)
	//}
	//fmt.Printf("ss is %v", ss)

}

//func function() {
//	// TODO 具体逻辑
//	fmt.Println("11111")
//	// 每5分钟执行一次
//	time.AfterFunc(5*time.Minute, function)
//}
