package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/routes"
	"redisData/setting"
)

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

	fmt.Println("success")
	//初始化routes
	r := routes.SetUp()
	r.Run(fmt.Sprintf(":%d",viper.GetInt("port")))

	//宕机处理
	defer func() {
		recover()
	}()
	//自动触发接口

}
