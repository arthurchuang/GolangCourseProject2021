package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
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
			pageUrl := fmt.Sprintf("%s%s", baseUrl, addr)
			if err := processPage(pageUrl); err != nil {
				fmt.Printf("Error while processing page (%s) : %e", pageUrl, err)
			}
		}
	})
}

func processPage(url string) error {
	fmt.Printf("page url: %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status code error: %d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	doc.Find(".hot-recommend-item.line").Each(func(i int, selection *goquery.Selection) {
		nameAnchor := selection.Find(".commodity-desc").Find("a")
		name, found := nameAnchor.Attr("title")
		if !found {
			name = "not found"
		}

		link, found := nameAnchor.Attr("href")
		if !found {
			link = "not found"
		}

		img := selection.Find(".gtm-product-alink").Find("img")
		imgLink, found := img.Attr("src")
		if !found {
			imgLink = "not found"
		}

		price := selection.Find(".current-price").Find("em").Text()

		saveEntry(name, link, imgLink, price)
	})
	return nil
}

func saveEntry(name string, link string, imgLink string, price string) {
	fmt.Printf("product name: %s\n", name)
	fmt.Printf("product link: %s%s\n", baseUrl, link)
	fmt.Printf("product image: %s\n", imgLink)
	fmt.Printf("product price: %s", price)
	fmt.Printf("\n\n")

	// TODO : save to DB
}
