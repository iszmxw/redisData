package model

type ApiKlineParam struct {
	Size   int64  `json:"size" form:"size"`
	Time   int64  `json:"time" form:"time"`
	Period string `json:"period" form:"period"`
}

type WebSocketKlineParam struct {
	Period string `json:"period" form:"period"` //根据时间返回不同的k线图
	Time   int64  `json:"time" form:"time"`     //返回给客户端的时间
}
