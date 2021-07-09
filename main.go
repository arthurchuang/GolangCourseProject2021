package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	baseURL := "https://www.carrefour.com.tw/"
	response, err := http.Get(baseURL)
	if err != nil {
		log.Fatal("Unable to parse from the baseURL: ", err)
	}
	body, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Fatal("Failed to read from HTML's body: ", errRead)
	}
	fmt.Println(string(body))
	fmt.Println("harry hello")

	response.Body.Close()
}
