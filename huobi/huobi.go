package huobi

import (
	"github.com/bitly/go-simplejson"
	"redisData/huobi/client"
	"redisData/huobi/market"
)

type JSON = simplejson.Json
type ParamsData = client.ParamData
type Market = market.Market
type Listener = market.Listener
type Client = client.Client

/// 创建WebSocket版Market客户端
func NewMarket() (*market.Market, error) {
	return market.NewMarket()
}

/// 创建RESTFul客户端
func NewClient(accessKeyId, accessKeySecret string) (*client.Client, error) {
	return client.NewClient(client.Endpoint, accessKeyId, accessKeySecret)
}
