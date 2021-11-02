package routes

import (
	"github.com/gin-gonic/gin"
	"redisData/controller"
)

func SetUp() *gin.Engine {
	r := gin.Default()

	//测试
	//r.GET("/hello", controller.Hello)
	//r.GET("/test", controller.Test)
	//r.GET("/ping", controller.Ping)

	//查询，查询redis上的数据，返回给前端
	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
	r.GET("/getRedisData", controller.GetRedisData)
	//r.GET("/autoGetRedisData/:period/:size/:symbol", controller.AutoGetRedisData)
	return r

}
