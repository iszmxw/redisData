package huobi

import (
	"encoding/json"
	"fmt"
	"github.com/leizongmin/huobiapi"
	"redisData/dao/mysql"
	"redisData/dao/redis"
)

type SubData struct {
	Ch    string `json:"ch"`
	Ts    int64  `json:"ts"`
	*Tick `json:"tick"`
}

type Tick struct {
	ID     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int64   `json:"count"`
}

func NewSubscribe(symbol string, period string) {
	//参数校验
	if period == "" {
		period = "1min"
	}
	// 创建客户端实例
	market, err := huobiapi.NewMarket()
	if err != nil {
		panic(err)
	}
	// 订阅主题
	market.Subscribe(fmt.Sprintf("market.%s.kline.%s", symbol, period), func(topic string, hjson *huobiapi.JSON) {
		// 收到数据更新时回调
		//fmt.Println(topic, json)

		//redis.CreateOrChangeKline(topic, *json)
		//fmt.Println(json)
		//utils.JSONToMap(string(json.MarshalJSON()))
		//jsonData是订阅后返回的信息 通过MarshalJSON将数据转化成String
		jsonData, _ := hjson.MarshalJSON()
		//mapData := utils.JSONToMap(string(jsonData))
		subData := &SubData{}
		if err := json.Unmarshal(jsonData, subData); err != nil {
			fmt.Printf("json.Unmarshal subData fail err:%v", err)
		}
		//通过数据库得到 自有币位数
		decimalscale, err := mysql.GetDecimalScaleBySymbols(symbol)
		if err != nil {
			fmt.Printf("mysql.GetDecimalScaleBySymbols fail err:%v", err)
			return
		}
		//对数据和自有币位数进行运算，返回修改后的数据
		if decimalscale.Value > 0 {
			subData.Amount = subData.Amount * float64(decimalscale.Value) * 0.01
			subData.Open = subData.Open * float64(decimalscale.Value) * 0.01
			subData.Close = subData.Close * float64(decimalscale.Value) * 0.01
			subData.Low = subData.Low * float64(decimalscale.Value) * 0.01
			subData.High = subData.High * float64(decimalscale.Value) * 0.01
			subData.Vol = subData.Vol * float64(decimalscale.Value) * 0.01
		}
		if decimalscale.Value < 0 {
			subData.Amount = subData.Amount / float64(decimalscale.Value) * 0.01
			subData.Open = subData.Open * float64(decimalscale.Value) * 0.01
			subData.Close = subData.Close * float64(decimalscale.Value) * 0.01
			subData.Low = subData.Low * float64(decimalscale.Value) * 0.01
			subData.High = subData.High * float64(decimalscale.Value) * 0.01
			subData.Vol = subData.Vol * float64(decimalscale.Value) * 0.01
		}
		//把修改后的对象反序列化，存进redis
		redisData, err := json.Marshal(subData)
		if err != nil {
			fmt.Printf("json.Marshal(subData) fail err:%v", err)
		}
		//根据推送返回的数据，以字符串的形式存入reids
		redis.CreateOrChangeKline(topic, string(redisData))

	})

	market.Loop()
}
