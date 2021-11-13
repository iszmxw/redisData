
package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/leizongmin/huobiapi"
	"net/http"
	"redisData/huobi"
	"redisData/logic"
	"redisData/model"
	"redisData/utils"
	"strings"
	"sync"
	"time"
)

//声明一个线程安全的map,存放ws和user
var users sync.Map


//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader1 = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsConn1 声明并发安全的ws
type WsConn1 struct {
	*websocket.Conn
	Mux sync.RWMutex
}

// UserInfo 看这个用户订阅了什么
type UserInfo struct {
	Uid     string `json:"uid"`
	Sub    []string `json:"sub_topic"`
}

func GetRedisData2(c *gin.Context)  {
	//每个用户连接,就new一个 ws
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err !=nil{
		fmt.Printf("upGrader.Upgrade err%v",err)
		return
	}
	defer ws.Close() //返回前关闭
	var user UserInfo
	user.Uid = utils.GetGenerateId()
	wsConn := &WsConn{
		ws,
		sync.RWMutex{},
	}

	//每一个ws 对应一个market
	//连接一个market
	market, err := huobi.NewMarket()
	if err != nil {
		fmt.Printf("huobi.NewMarket() fail %v",err)
	}
	//把这个websocket指针和user对应的hash储存起来
	users.Store(ws,user)

	//读取客户端信息
	for {
		//读取ws中的数据
		wsConn.Mux.Lock()
		mt, message, err := wsConn.Conn.ReadMessage()
		wsConn.Mux.Unlock()
		if err != nil {
			fmt.Println(err)
			break
		}

		//把用户传进来的消息进行处理 msg样式 "market.btcusdt.kline.1min"
		msg := string(message)
		//-------------
		fmt.Println(msg)
		//当请求数据中含有1min或1step这些为已经缓存数据,直接去redis拿
		if strings.Contains(msg,"1min")||strings.Contains(msg,"step1"){
			go func() {
				for  {
					data, err := logic.GetDataByKey(msg)
					if err !=nil{
						//如果redis数据获取或者start接口没有被调用，就要重新缓存
						if err == redis.Nil{
							wsConn.Mux.Lock()
							err = wsConn.Conn.WriteMessage(mt, []byte("数据已过期，准备开始缓存"))
							wsConn.Mux.Unlock()
							if err != nil {
								fmt.Printf("wsConn.Conn.WriteMessage fail %v",err)
								return
							}
							fmt.Printf("logic.GetDataByKey fail %v",err)
							//5s 订阅一次，避免newMarket报错
							//订阅k线图的数据
							err := logic.StartSetKlineData()
							if err != nil {
								fmt.Printf("logic.StartSetKlineData fail err%v",err)
								return
							}
							time.Sleep(2*time.Second)
							//订阅行情的数据
							err = logic.StartSetQuotation()
							if err != nil {
								fmt.Printf("logic.StartSetQuotation fail err%v",err)
								return
							}
							time.Sleep(10 * time.Second)
						}
					}
					//把读到数据，序列化后返回
					websocketData := utils.Strval(data)
					wsConn.Mux.Lock()
					err = wsConn.Conn.WriteMessage(mt, []byte(websocketData))
					wsConn.Mux.Unlock()
					if err != nil {
						fmt.Println(err)
						ws.Close()

					}
					//每2s推送一次
					time.Sleep(time.Second * 2)
				}
			}()
			return
		}
//第二部分逻辑  输入参数为"market.btcusdt.kline.5min" ，等数据库不存在的数据，直接转发60秒后取消订阅，刷新后重新订阅（先不做看性能如何）
		//如果请求的是market.ethbtc.kline.5min,订阅这条信息，然后再返回
		//msg 用户输入的数据byte转string
		//直接拿msg去订阅
		market.Subscribe(msg, func(topic string, hjson *huobiapi.JSON) {
			// 收到数据更新时回调
			fmt.Println(topic, hjson)
			jsondata, err := hjson.MarshalJSON()
			if err != nil {
				fmt.Printf("hjson.MarshalJSON fail err%v",err)
				return
			}
			//把jsondata反序列化后进行，自由币判断运算
			klineData := model.KlineData{}
			err = json.Unmarshal(jsondata, &klineData)
			if err != nil {
				fmt.Printf("json.Unmarshal %v",err)
				return
			}
			//自由币换算
			tranData := logic.TranDecimalScale2(msg,klineData)
			//结构体序列化后返回
			data, err := json.Marshal(tranData)
			if err != nil {
				fmt.Printf("json.Marshal(tranData) fail %v",err)
				return
			}
			//返回数据给用户
			wsConn.Mux.Lock()
			err = wsConn.Conn.WriteMessage(mt, data)
			wsConn.Mux.Unlock()
			if err != nil {
				fmt.Println(err)
				ws.Close()

			}

		})


	}
	//关闭前处理下订阅信息
	//定义回调函数
	h := func(code int,text string ) error {
		//最好加一个取消订阅
		err := market.Close()
		if err != nil {
			fmt.Printf("market.Close() fail %v",err)
		}
		return nil
	}
	ws.SetCloseHandler(h)




}

