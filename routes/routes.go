package routes

import (
	"github.com/gin-gonic/gin"
	"redisData/controller"
)

func SetUp() *gin.Engine {
	r := gin.Default()

	//查询，查询redis上的数据，返回给前端
	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
	r.GET("/getRedisData", controller.GetRedisData)
	r.GET("/start", controller.StartController)
	r.GET("/quotation", controller.QuotationController)
	return r

}
