package main

import (
	"GoCrawl/crawl"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const (
	storeUrl = "https://online.carrefour.com.tw/tw/"
	baseUrl  = "https://online.carrefour.com.tw"
)

func main() {
	numWorkers := flag.Int("numWorkers", 3, "Use this flag to set the number of workers (default to 3 if not specified).")
	flag.Parse()
	fmt.Printf("Number of workers: %d\n", *numWorkers)

	doc, err := crawl.GetUrlDocument(storeUrl)
	if err != nil {
		log.Fatalf("Failed to read from store url (%s): %e", storeUrl, err)
	}

	start := time.Now()

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

	elapsed := time.Since(start)
	fmt.Printf("\n\nTook %s to process all categories with %d workers.", elapsed, *numWorkers)
}

func processPage(url string) error {
	fmt.Printf("page url: %s\n", url)
	doc, err := crawl.GetUrlDocument(url)
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
	fmt.Printf("product price: %s\n\n", price)

	// TODO : save to DB
}
