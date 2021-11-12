package main

import (
	"fmt"
	"github.com/leizongmin/huobiapi"
)

//var upGrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}

//type WsConn struct {
//	*websocket.Conn
//	Mux sync.RWMutex
//}

func main() {
	// 创建客户端实例
	market, err := huobiapi.NewMarket()
	if err != nil {
		panic(err)
	}
	// 订阅主题
	market.Subscribe("market.ethbtc.kline.1min", func(topic string, json *huobiapi.JSON) {
		// 收到数据更新时回调
		fmt.Println(topic, json)
	})
	// 请求数据
	json, err := market.Request("market.btcusdt.detail")
	if err != nil {
		panic(err)
	} else {
		fmt.Println(json)
	}
	// 进入阻塞等待，这样不会导致进程退出
	market.Loop()
}


