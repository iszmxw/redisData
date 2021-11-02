package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"redisData/logic"
	"redisData/utils"
	"time"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//func AutoGetRedisData(c *gin.Context) {
//	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
//	//校验参数，不写直接给默认值
//
//	//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器
//
//	//返回数据
//}

func GetRedisData(c *gin.Context) {

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close() //返回前关闭

	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//自定义修改 1.获取参数 2.调用逻辑 3.返回数据
		data, err := logic.GetDataByKey(string(message))

		websocketData := utils.Strval(data)
		//时间参数
		//var t D
		//times := c.Param("times")
		//t, err := strconv.Atoi(times)
		if err != nil {

			fmt.Printf("字符串转换int类型失败,err is %v", err)
		}
		//写入ws数据
		go func() {
			for {
				err = ws.WriteMessage(mt, []byte(websocketData))
				if err != nil {
					return
				}
				time.Sleep(time.Second * 1)
			}

		}()

	}

	//校验参数，不写直接给默认值
	//symbol := c.Param("symbol")
	//if symbol == "" {
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "symbol param is require",
	//	})
	//	return
	//}
	//判断key是否存在

	//逻辑处理
	//data, err := logic.GetDataByKey(symbol)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": err,
	//	})
	//	return
	//}
	////返回数据
	//c.JSON(http.StatusOK, gin.H{
	//	"redisData": data,
	//})

}
