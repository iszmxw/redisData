package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type websocketClientManager struct {
	conn        *websocket.Conn
	addr        *string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

// 构造函数
func NewWsClientManager(addrIp string, timeout int) *websocketClientManager {
	addrString := addrIp
	var sendChan = make(chan string, 10)
	var recvChan = make(chan string, 10)
	var conn *websocket.Conn
	return &websocketClientManager{
		addr:        &addrString,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
	}
}

// 链接服务端
func (wsc *websocketClientManager) dail() {
	var err error
	u := *wsc.addr
	log.Printf("connecting to %s", u)
	wsc.conn, _, err = websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		fmt.Println(err)
		return

	}
	wsc.isAlive = true
	log.Printf("connecting to %s 链接成功！！！", u)

}

// 发送消息
func (wsc *websocketClientManager) sendMsgThread() {
	go func() {
		for {
			msg := <-wsc.sendMsgChan
			err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				continue
			}
		}
	}()
}

// 读取消息
func (wsc *websocketClientManager) readMsgThread() {
	go func() {
		for {
			if wsc.conn != nil {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				log.Printf("recv: %s", message)
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)
			}

		}
	}()
}

// 开启服务并重连
func (wsc *websocketClientManager) start() {
	for {
		if wsc.isAlive == false {
			wsc.dail()
			wsc.sendMsgThread()
			wsc.readMsgThread()
		}
		time.Sleep(time.Second * time.Duration(wsc.timeout))
	}
}

func main() {
	wsc := NewWsClientManager("wss://api.huobi.pro/ws", 10)
	wsc.start()
	//1.发送ping加时间戳
	//2.接受数据
	//3.判断数据类型
	//4.如果返回的是ping,就发送pong
	//5.发送订阅信息
	//6.
	wsc.readMsgThread()
	var w1 sync.WaitGroup
	w1.Add(1)
	w1.Wait()
}
