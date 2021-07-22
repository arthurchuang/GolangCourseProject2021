package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	storeUrl = "https://online.carrefour.com.tw/tw/"
	baseUrl  = "https://online.carrefour.com.tw"
)

type productEntry struct {
	category, name, link, imgLink, price string
}

func main() {
	t1 := time.Now()
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

	// Declare a slice of strings to contain each category collected from query
	var category []string
	doc.Find(".top1.left-item").Each(func(i int, selection *goquery.Selection) {
		anchor := selection.Find("a")
		addr, found := anchor.Attr("href")
		if found {
			decodedValue, errValue := url.QueryUnescape(addr)
			if errValue != nil {
				log.Fatal(errValue)
			}
			category = append(category, decodedValue)
		}
	})

	// Define the number of jobs as the number of categories
	numJobs := len(category)
	jobs := make(chan string, numJobs)
	results := make(chan []productEntry, numJobs)

	for w := 1; w <= 5; w++ {
		go worker(w, jobs, results)
	}

	j := 0
	for range category {
		jobs <- category[j]
		j++
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		items := <-results
		for _, v := range items {
			saveEntry(v.category, v.name, v.link, v.imgLink, v.price)
		}
	}

	t2 := time.Now()

	fmt.Println(t1)
	fmt.Println(t2)

	diff := t2.Sub(t1)
	fmt.Println(diff)
}

func worker(id int, jobs <-chan string, results chan<- []productEntry) {
	for j := range jobs {
		// fmt.Println("worker", id, "started  job", j)
		pageUrl := fmt.Sprintf("%s%s", baseUrl, j)
		out, err := processPage(pageUrl)
		if err != nil {
			fmt.Printf("Error while processing page (%s) : %e", pageUrl, err)
		}
		// fmt.Println("worker", id, "finished job", j)
		results <- out
	}
}

func processPage(url string) ([]productEntry, error) {
	fmt.Printf("page url: %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status code error: %d %s", res.StatusCode, res.Status))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var productList []productEntry
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

		productList = append(productList, productEntry{url, name, link, imgLink, price})
	})
	return productList, nil
}

func saveEntry(category string, name string, link string, imgLink string, price string) {
	fmt.Printf("category url: %s\n", category)
	fmt.Printf("product name: %s\n", name)
	fmt.Printf("product link: %s%s\n", baseUrl, link)
	fmt.Printf("product image: %s\n", imgLink)
	fmt.Printf("product price: %s", price)
	fmt.Printf("\n\n")

	// TODO : save to DB
}
