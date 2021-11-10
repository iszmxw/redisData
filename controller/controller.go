package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

// StartController 每次先手动启动一下
func StartController(c *gin.Context) {

	//启动获取k线图数据
	if err := logic.StartSetKlineData(); err != nil {
		fmt.Printf("logic.StartSetKlineData() fail err:%v", err)
	}
	//启动获取行情数据
	if err := logic.StartSetQuotation(); err != nil {
		fmt.Printf("logic.StartSetQuotation() fail err:%v", err)
	}
	fmt.Println("huobiService Start success")
	c.JSON(http.StatusOK, gin.H{

		"msg": "Start success",
	})
}

// GetRedisData websocket请求,根据发送的内容返回键值对
func GetRedisData(c *gin.Context) {

	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close() //返回前关闭
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//对数据进行切割，读取参数
		//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
		strMsg := string(message)
		//打印请求参数
		fmt.Println(strMsg)

		//写入ws数据
		go func() {
			for {
				data, err := logic.GetDataByKey(strMsg)
				//修改，当拿不到key重新订阅，10秒订阅一次
				if err == redis.Nil {
					logic.StartSetKlineData()
					time.Sleep(10 * time.Second)
				}
				websocketData := utils.Strval(data)
				err = ws.WriteMessage(mt, []byte(websocketData))
				if err != nil {
					return
				}
				time.Sleep(time.Second * 1)
			}

		}()

	}

}

// QuotationController 请求行情数据接口
func QuotationController(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close() //返回前关闭
	for {
		//读取ws中的数据，数据是"market.btcusdt.depth.step1"类型
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		strMsg := string(message)
		//打印请求参数
		fmt.Println(strMsg)
		//分割
		//resultList := utils.Split(strMsg, ".")

		go func() {
			for {
				data, err := logic.GetDataByKey(strMsg)
				//修改，当拿不到key重新订阅，10秒订阅一次
				if err == redis.Nil {
					logic.StartSetQuotation()
					time.Sleep(10 * time.Second)
				}
				websocketData := utils.Strval(data)
				err = ws.WriteMessage(mt, []byte(websocketData))
				if err != nil {
					return
				}
				time.Sleep(time.Second * 1)
			}

		}()
	}
}
