package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const (
	storeUrl = "https://online.carrefour.com.tw/tw/"
	baseUrl  = "https://online.carrefour.com.tw"
)

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

	doc.Find(".top1.left-item").Each(func(i int, selection *goquery.Selection) {
		anchor := selection.Find("a")
		addr, found := anchor.Attr("href")
		if found {
			pageLink := fmt.Sprintf("%s%s", baseUrl, addr)
			fmt.Printf("page link: %s\n", pageLink)
			// TODO : process page
		}
	})
}
