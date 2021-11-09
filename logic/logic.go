package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"redisData/dao/mysql"
	"redisData/dao/redis"
	"redisData/huobi"
	"redisData/utils"
	"strings"
)

var (
	ErrorUnmarshalFail = errors.New("UnmarshalFail")
)

// StartSetKlineData main.go时，默认缓存16种symbol,1min的数据
func StartSetKlineData(period string) error {
	//通过访问mysql获取切片
	symbol, err := mysql.GetAllSymbol()
	if err != nil {
		fmt.Printf("mysql.GetAllSymbol fail %v", err)
		return err
	}
	ss := make([]string, 0, len(*symbol))
	for _, value := range *symbol {
		ss = append(ss, value.Name)
	}
	//根据symbol切片长度起goroutine
	//1.遍历mysql中的symbol,NewSubscribe中有存入redis的方法
	for i := 0; i < len(ss); i++ {
		go huobi.NewSubscribe(ss[i], period)
	}
	return nil
}

// StartSetQuotation 自动获取行情数据
func StartSetQuotation() error {
	//通过访问mysql获取切片
	symbol, err := mysql.GetAllSymbol()
	if err != nil {
		fmt.Printf("mysql.GetAllSymbol fail %v", err)
		return err
	}
	ss := make([]string, 0, len(*symbol))
	for _, value := range *symbol {
		ss = append(ss, value.Name)
	}
	//根据symbol切片长度起goroutine
	//1.遍历mysql中的symbol,NewQuotation中有存入redis的数据中
	for i := 0; i < len(ss); i++ {
		go huobi.NewQuotation(ss[i])
	}
	return nil
}

//SubscribeOneKline 客户端连接socket后，向服务端发送数据
func SubscribeOneKline(symbol string, period string) {
	go huobi.NewSubscribe(symbol, period)
	//if <-ch != 1 {
	//	market.Unsubscribe(fmt.Sprintf("market.%s.kline.%s", symbol, period))
	//	market.Close()
	//	return
	//}
}

// AutoGetRedisData 通过HTTP请求自动存取symbol数据到redis
func AutoGetRedisData(size int64, period string) error {
	if size == 0 {
		size = 300
	}
	if period == "" {
		period = "1min"
	}
	//通过访问mysql获取切片
	symbol, err := mysql.GetAllSymbol()
	if err != nil {
		fmt.Printf("mysql.GetAllSymbol fail %v", err)
		return err
	}
	ss := make([]string, 0, len(*symbol))
	for _, value := range *symbol {
		ss = append(ss, value.Name)
	}
	//fmt.Printf("ss is %v", ss)
	//fmt.Printf("ss is %T", ss)

	//把这部分修改成websocket

	//传入切片，拼接url参数发起请求，把数据存进redis
	for i := 0; i < len(ss); i++ {
		url := fmt.Sprintf("https://api.huobi.pro/market/history/kline?period=%s&size=%d&symbol=%s", period, size, ss[i])
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("get api fail err is %v", err)
			return err
		}
		body, _ := ioutil.ReadAll(response.Body)
		data := string(body)

		//把数据写进redis
		fmt.Println("redis开始写数据")
		redis.CreateOrChangeKline(ss[i], data)
		fmt.Println("redis结束写数据")

	}

	//for _, v := range ss {
	//	url := fmt.Sprintf("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=%s", v)
	//	response, err := http.Get(url)
	//	if err != nil {
	//		log.Fatalf("get api fail err is %v", err)
	//		return err
	//	}
	//	body, _ := ioutil.ReadAll(response.Body)
	//	data := string(body)
	//
	//	//把数据写进redis
	//	redis.CreateOrChangeKline(v, data)
	//	return nil
	//}
	return nil

}

func GetDataByKey(key string) (interface{}, error) {
	//根据key获取值
	kline, err := redis.GetKline(key)
	if err != nil {
		return nil, err
	}
	//将对应key中的value值，将string转化成json后返回
	data := []byte(kline)
	var i interface{}
	//3.解析
	if err := json.Unmarshal(data, &i); err != nil {
		fmt.Println(err)
		return nil, ErrorUnmarshalFail
	}
	return i, nil
}

//CheckDataType 区分ping请求和订阅请求
func CheckDataType(str string) (dataType int, str2 string) {
	if strings.Contains(str, "ping") {
		rest := utils.Split(str, ":")
		str2 = rest[1]
		return 1, str2
	}
	if strings.Contains(str, "pong") {
		return 2, "success"
	}
	if strings.Contains(str, "kline") {
		return 3, str
	}
	return 0, "other"
}
