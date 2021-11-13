package model


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
