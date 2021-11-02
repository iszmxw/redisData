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
)

var (
	ErrorUnmarshalFail = errors.New("UnmarshalFail")
)

func AutoGetRedisData() error {

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
	fmt.Printf("ss is %v", ss)
	fmt.Printf("ss is %T", ss)

	//传入切片，拼接url参数发起请求，把数据存进redis
	for i := 0; i < len(ss); i++ {
		url := fmt.Sprintf("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=%s", ss[i])
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("get api fail err is %v", err)
			return err
		}
		body, _ := ioutil.ReadAll(response.Body)
		data := string(body)

		//把数据写进redis
		redis.CreateOrChangeKline(ss[i], data)

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
