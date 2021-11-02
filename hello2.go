package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	//url := "127.0.0.1:8081/version"
	response, err := http.Get("https://api.huobi.pro/market/history/kline?period=1min&size=1&symbol=btcusdt")

	//req, err := http.NewRequest("GET", url, nil)

	//client := &http.Client{}

	//resp, err := client.Do(req)

	if err != nil {

		panic(err)

	}

	defer response.Body.Close()

	//fmt.Println("response Status:", resp.Status)

	//fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println("response Body:", string(body))

}
