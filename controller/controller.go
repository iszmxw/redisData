package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"redisData/dao/redis"
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

//func QueryKlineData(size int64,period string) {
//	//response, err := http.Get("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt")
//
//	if err := logic.AutoGetRedisData(size,period); err != nil {
//		fmt.Printf("AutoGetRedisData is fail %v", err)
//		return
//	}
//	time.AfterFunc(30*time.Second, QueryKlineData(size,period))
//	log.Println("开始获取交易对数据")
//}

//func AutoGetRedisData(c *gin.Context) {
//	//请求的 url https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt
//	//校验参数，不写直接给默认值
//
//	//逻辑处理 每30s发起请求 起10多个goroutine同时发起请求 通过数据拼接发起请求 拿到响应数据 存入redis同时发送到服务器
//
//	//返回数据
//}

func GetRedisData(c *gin.Context) {

	//使用websocket后废弃
	//------------------------------
	//获取参数URL参数
	//p := &model.WebSocketKlineParam{
	//	Time:   1,
	//	Period: "1min",
	//}
	//if err := c.ShouldBindQuery(p); err != nil {
	//	fmt.Printf("ShouldBindQuery fail err:%v", err)
	//	return
	//}
	//------------------------------

	//把获取数据的逻辑放在这里,根据参数把数据存进redis 写入logic层
	//通过websocket访问
	if err := logic.StartSetRedisData(""); err != nil {
		fmt.Printf("logic.StartSetRedisData() fail err:%v", err)
	}

	//通过http访问
	//go func() {
	//	for {
	//		if err := logic.AutoGetRedisData(p.Size, p.Period); err != nil {
	//			fmt.Printf("AutoGetRedisData is fail %v", err)
	//			return
	//		}
	//		time.Sleep(30 * time.Second)
	//	}
	//
	//}()

	//ch := make(chan *model.ApiKlineParam)
	//go QueryKlineData(ch)

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
		//对数据进行切割，读取参数
		//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
		strMsg := string(message)
		//判断下key是否缓存过，否的话提交缓存
		ret := redis.ExistKey(strMsg)
		if ret != true {
			resultList := utils.Split(strMsg, ".")
			logic.SubscribeOneKline(resultList[1], resultList[3])
		}
		//自定义修改 1.获取参数 2.调用逻辑 3.返回数据
		//-------------------------------------------------------
		//通过key获取数据后返回
		//data, err := logic.GetDataByKey(strMsg)
		//websocketData := utils.Strval(data)
		//-------------------------------------------------------
		//时间参数
		//times := p.Time
		//n := utils.GetSleepTime(times)
		if err != nil {
			fmt.Printf("字符串转换int类型失败,err is %v", err)
		}
		//写入ws数据
		go func() {
			for {
				data, err := logic.GetDataByKey(strMsg)
				websocketData := utils.Strval(data)
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
