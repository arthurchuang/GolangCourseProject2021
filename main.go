package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const storeUrl = "https://online.carrefour.com.tw/tw/"

func main() {
	res, err := http.Get(storeUrl)
	if err != nil {
		log.Fatalf("Error getting store website: %e", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// TODO : locate actual things of interest
	doc.Find(".menu-list").Each(func(i int, selection *goquery.Selection) {
		fmt.Printf("index %d, content %v", i, selection)
	})
}
