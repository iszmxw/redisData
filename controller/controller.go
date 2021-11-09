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

// StartController 每次先手动启动一下
func StartController(c *gin.Context) {

	//启动获取k线图数据
	if err := logic.StartSetRedisData(""); err != nil {
		fmt.Printf("logic.StartSetRedisData() fail err:%v", err)
	}
	//启动获取行情数据
	if err := logic.StartSetQuotation(); err != nil {
		fmt.Printf("logic.StartSetQuotation() fail err:%v", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Start success",
	})
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

// GetRedisData websocket请求,根据发送的内容返回键值对
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

	//
	//if err := logic.StartSetRedisData(""); err != nil {
	//	fmt.Printf("logic.StartSetRedisData() fail err:%v", err)
	//}
	//

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
		fmt.Println(strMsg)
		//判断下key是否缓存过，否的话提交缓存
		//ret := redis.ExistKey(strMsg)
		//if ret != true {
		//	resultList := utils.Split(strMsg, ".")
		//	logic.SubscribeOneKline(resultList[1], resultList[3])
		//}
		//自定义修改 1.获取参数 2.调用逻辑 3.返回数据
		//-------------------------------------------------------
		//通过key获取数据后返回
		//data, err := logic.GetDataByKey(strMsg)
		//websocketData := utils.Strval(data)
		//-------------------------------------------------------
		//时间参数
		//times := p.Time
		//n := utils.GetSleepTime(times)
		//if err != nil {
		//	fmt.Printf("字符串转换int类型失败,err is %v", err)
		//}
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
		//resultList := utils.Split(strMsg, ".")
		//ch := make(chan int)
		//ch <- 1
		//logic.SubscribeOneKline(resultList[1], resultList[3], ch)
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
	//ws.PingHandler()
}

//func GetRedisData2(c *gin.Context) {
//	//客户端部分代码
//	//把获取数据的逻辑放在这里,根据参数把数据存进redis 写入logic层
//	//通过websocket访问
//	if err := logic.StartSetRedisData(""); err != nil {
//		fmt.Printf("logic.StartSetRedisData() fail err:%v", err)
//	}
//
//	//服务端部分代码
//	//升级get请求为webSocket协议
//	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
//	if err != nil {
//		return
//	}
//	defer ws.Close() //返回前关闭
//	ch := make(chan int)
//	ch <- 1
//	for {
//
//		//读取ws中的数据
//		mt, message, err := ws.ReadMessage()
//		if len(message) == 0 {
//			pingData := fmt.Sprintf("{\"ping\"}:%s", utils.Strval(time.Now().Unix()))
//			ws.WriteMessage(mt, []byte(pingData))
//		}
//		if err != nil {
//			break
//		}
//		strMsg := string(message)
//		//对数据参数进行分类 1.是ping 2.是pong 3.是market.ethbtc.kline.1min
//		t, s := logic.CheckDataType(strMsg)
//		if t == 1 {
//			rtData := fmt.Sprintf("{\"ping\"}:%s", s)
//			ws.WriteMessage(mt, []byte(rtData))
//		}
//		if t == 2 {
//			c.JSON(http.StatusOK, s)
//		}
//		if t == 3 {
//			//对数据进行切割，读取参数
//			//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
//			//strMsg := string(message)
//			//判断下key是否缓存过，否的话提交缓存
//			ret := redis.ExistKey(s)
//			if ret != true {
//				resultList := utils.Split(s, ".")
//				logic.SubscribeOneKline(resultList[1], resultList[3], ch)
//			}
//		}
//		if err != nil {
//			fmt.Printf("字符串转换int类型失败,err is %v", err)
//		}
//
//		//写的协程,写入ws数据
//		go func() {
//			for {
//				data, err := logic.GetDataByKey(strMsg)
//				websocketData := utils.Strval(data)
//				err = ws.WriteMessage(mt, []byte(websocketData))
//				if err != nil {
//					return
//				}
//				time.Sleep(time.Second * 1)
//			}
//		}()
//		//一直读
//		go func() {
//			_, msg, err := ws.ReadMessage()
//			time.Sleep(10 * time.Second)
//
//			if len(msg) == 0 {
//				_, msg, _ = ws.ReadMessage()
//				if len(msg) == 0 {
//					ws.Close()
//					ch <- 0
//					return
//				}
//				//pingData := fmt.Sprintf("{\"ping\"}:%s", utils.Strval(time.Now().Unix()))
//				//ws.WriteMessage(mt, []byte(pingData))
//			}
//			if err != nil {
//				return
//			}
//		}()
//	}
//	//go func() {
//	//	for {
//	//		pingData := fmt.Sprintf("{\"ping\"}:%s", utils.Strval(time.Now().Unix()))
//	//		err := ws.WriteMessage(mt, []byte(pingData))
//	//		if err != nil {
//	//			fmt.Printf("ws.WriteMessage fail err:v%", err)
//	//			return
//	//		}
//	//		time.Sleep(10 * time.Second)
//	//	}
//	//}()
//
//	//for {
//	//
//	//	go func() {
//	//		_, msg, err := ws.ReadMessage()
//	//		if err != nil {
//	//			return
//	//		}
//	//		if strings.Contains(string(msg), "pong") {
//	//			ws.Close()
//	//			return
//	//		}
//	//	}()
//	//	time.Sleep(10 * time.Second)
//	//
//
//}

// QuotationController 请求行情数据接口
func QuotationController(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close() //返回前关闭
	for {
		//读取ws中的数据，数据是"market.btcusdt.depth.step1"
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		strMsg := string(message)
		//分割
		//resultList := utils.Split(strMsg, ".")

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
}
