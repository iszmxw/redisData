package huobi

import (
	"encoding/json"
	"fmt"
	"github.com/leizongmin/huobiapi"
	"github.com/leizongmin/huobiapi/market"
	"net/http"
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

type Quotation struct {
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	*Ticks `json:"tick"`
}
type Ticks struct {
	Bids    [][]float64 `json:"bids"`
	Asks    [][]float64 `json:"asks"`
	Version int64       `json:"version"`
	Ts      int64       `json:"ts"`
}

// NewSubscribe 新订阅
func NewSubscribe(symbol string, period string) *market.Market {
	//fmt.Printf("market.%s.kline.%s", symbol, period)
	//参数校验
	if period == "" {
		period = "1min"
	}
	// 创建客户端实例
	market, err := huobiapi.NewMarket()
	if err != nil {
		fmt.Println(err)
		err = market.ReConnect()
		http.Get("localhost:8887/start")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
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
			decimalscale.Value = decimalscale.Value * -1
			subData.Amount = subData.Amount / float64(decimalscale.Value) * 0.01
			subData.Open = subData.Open / float64(decimalscale.Value) * 0.01
			subData.Close = subData.Close / float64(decimalscale.Value) * 0.01
			subData.Low = subData.Low / float64(decimalscale.Value) * 0.01
			subData.High = subData.High / float64(decimalscale.Value) * 0.01
			subData.Vol = subData.Vol / float64(decimalscale.Value) * 0.01
		}
		//把修改后的对象反序列化，存进redis
		redisData, err := json.Marshal(subData)
		if err != nil {
			fmt.Printf("json.Marshal(subData) fail err:%v", err)
		}
		//根据推送返回的数据，以字符串的形式存入reids
		//redis.CreateOrChangeKline(topic, string(redisData))
		redis.CreateOrChangeKline(fmt.Sprintf("\"%s\"", topic), string(redisData))
		//fmt.Println(string(redisData))

	})
	//c := <-ch
	//if c == 0 {
	//	market.Unsubscribe(fmt.Sprintf("market.%s.kline.%s"))
	//	market.Close()
	//}

	market.Loop()
	return market

}

func NewQuotation(symbol string) {
	// 创建客户端实例
	market, err := huobiapi.NewMarket()
	if err != nil {
		fmt.Println(err)
		err = market.ReConnect()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	//订阅行情信息
	//"market.btcusdt.depth.step0"
	market.Subscribe(fmt.Sprintf("market.%s.depth.step1", symbol), func(topic string, hjson *huobiapi.JSON) {
		// 收到数据更新时回调
		//fmt.Println(topic, hjson)
		jsonData, _ := hjson.MarshalJSON()
		//println(string(jsonData))

		data := new(Quotation)
		err := json.Unmarshal(jsonData, data)
		if err != nil {
			fmt.Println(err)
			return
		}
		//redis.CreateRedisData(fmt.Sprintf("\"%s\"", topic), string(jsonData))
		//fmt.Printf("%#v", data.Ticks.Bids[1][0])
		//根据自由币变量修改
		//通过数据库得到 自有币位数
		decimalscale, err := mysql.GetDecimalScaleBySymbols(symbol)
		if err != nil {
			fmt.Printf("mysql.GetDecimalScaleBySymbols fail err:%v", err)
			return
		}
		//对数据和自有币位数进行运算，返回修改后的数据
		if decimalscale.Value > 0 {
			for i := 0; i < len(data.Asks); i++ {
				data.Asks[i][0] = data.Asks[i][0] * float64(decimalscale.Value) * 0.01
			}
			for i := 0; i < len(data.Bids); i++ {
				data.Bids[i][0] = data.Bids[i][0] * float64(decimalscale.Value) * 0.01
			}

		}
		if decimalscale.Value < 0 {
			decimalscale.Value = decimalscale.Value * -1
			for i := 0; i < len(data.Asks); i++ {
				data.Asks[i][0] = data.Asks[i][0] / float64(decimalscale.Value) * 0.01
			}
			for i := 0; i < len(data.Bids); i++ {
				data.Bids[i][0] = data.Bids[i][0] / float64(decimalscale.Value) * 0.01
			}
		}
		//序列化，存进redis
		jsonData, err = json.Marshal(data)
		if err != nil {
			fmt.Printf("json.Marshal(data) fail,err%s", err)
		}
		redis.CreateRedisData(fmt.Sprintf("\"%s\"", topic), string(jsonData))

	})

}
