package model

type ApiKlineParam struct {
	Size   int64  `json:"size" `
	Symbol string `json:"symbol"`
	Period string `json:"period"`
}
