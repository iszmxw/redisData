package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "hello",
	})
}

func Test(c *gin.Context) {
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
}

func AutoGetRedisData(c *gin.Context) {
	//校验参数，不写直接给默认值

	//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器

	//返回数据
}

func GetRedisData(c *gin.Context) {
	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
	//校验参数，不写直接给默认值

	//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器

	//返回数据
}

