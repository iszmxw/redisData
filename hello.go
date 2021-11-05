package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"redisData/utils"
	"time"
)

type SubRequest struct {
	Sub string `json:"sub"`
	Id  string `json:"id"`
}

type SubResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Subbed string `json:"subbed"`
	Ts     int64  `json:"ts"`
}

type UpdateData struct {
	Ch    string `json:"ch"`
	Ts    string `json:"ts"`
	*Tick `json:"tick"`
}

type Tick struct {
	Id     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int64   `json:"count"`
}

func main() {
	url := "wss://api.huobi.pro/feed"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial", err)
	}
	defer conn.Close()

	timestamp := string(int32(time.Now().Unix()))
	msg1 := fmt.Sprintf("{\"ping\":%s}", timestamp)
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg1))
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	time.Sleep(2 * time.Second)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("err:%v", err)
	}
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Printf("WriteMessage is fail err:%v", err)
		return
	}

	s := SubRequest{
		Sub: "market.btcusdt.kline.1min",
		Id:  utils.GetGenerateId(),
	}
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("json.Marshal is fail err:%v", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("WriteMessage is fail err:%v", err)
		return
	}
	_, subResp, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("ReadMessage is fail err:%v", err)
		return
	}
	//var ss SubResponse
	//if err = json.Unmarshal(subResp, &ss); err != nil {
	//	fmt.Printf("json.Unmarshal is fail err:%v", err)
	//}
	t, _ := utils.ParseGzip(subResp)
	fmt.Println(string(t))
	for {
		_, updateData, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("ReadMessage is fail err:%v", err)
		}
		d, _ := utils.ParseGzip(updateData)
		//var u UpdateData
		//if err := json.Unmarshal(updateData, &u); err != nil {
		//	fmt.Printf("json.Unmarshal is fail err:%v", err)
		//}
		fmt.Println(string(d))
		time.Sleep(2 * time.Second)
	}
}
