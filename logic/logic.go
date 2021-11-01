package logic

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)



func AutoGetRedisData(symbols []string) error {
	//传入切片，拼接url参数发起请求，把数据存进redis
	//for _, v := range symbols{
	//
	//	url := fmt.Sprintf("api.huobi.pro/market/history/kline?period=1min&size=1&symbol=%s",v)
	//	response, err := http.Get(url)
	//	if err != nil{
	//		log.Fatalf("get api fail err is %v",err)
	//		return err
	//	}
	//	fmt.Println(response.Body)
		//把数据写进redis
		//redis.CreateKline(v,response)
	url := "https//:baidu.com"
	response, err := http.Get(url)

	if err != nil{
		log.Fatalf("get response fail err is %v",err)
				return err
	}
	body,_ := ioutil.ReadAll(response.Body)
	fmt.Println(body)
	return nil
}
