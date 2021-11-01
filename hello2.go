package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {


	url := "127.0.0.1:8081/version"

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {

		panic(err)

	}

	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)

	//fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))


}
