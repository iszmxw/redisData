package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"redisData/dao/redis"
)

var symbol = make([]string, 0)

//  QueryKlineData 每隔30s发送一次请求，请求全部类型的交易对
//  如果没有key就创建，存在key就更新

func QueryKlineData() {
	response, err := http.Get("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt")

}

func main() {
	//初始化redis
	err := redis.InitClient()
	if err != nil {
		fmt.Printf("init redis fail err:%v/n", err)
		return
	}

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		response, err := http.Get("https://baidu.com")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			//"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})

	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
	r.GET("/redisData/:period/:size/:symbol", func(c *gin.Context) {
		//校验参数，不写直接给默认值

		//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器

		//返回数据
	})
	//查询，查询redis上的数据，返回给前端
	r.GET("/getRedisData/:period/:size/:symbol", func(c *gin.Context) {
		//校验参数，不写直接给默认值

		//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器

		//返回数据
	})
	r.Run(":8887")
}
