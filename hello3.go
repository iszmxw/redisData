package main

import "fmt"

func main() {
	var symbols = []string{"btcusdt", "etcusdt", "xrpusdt", "adausdt", "ltcusdt", "xemusdt", "dashusdt", "xlmusdt", "ethbtc", "ethbtc", "eosbtc", "dotbtc", "dotbtc", "linketh", "adaeth", "xmreth"}
	for i := 0; i < len(symbols); i++ {

		url := fmt.Sprintf("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=%s", symbols[i])
		fmt.Println(i, url)
	}
}
