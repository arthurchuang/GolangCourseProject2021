package main

import (
	"GoCrawl/concurrent"
	"GoCrawl/crawl"
	"GoCrawl/database"
	"GoCrawl/model"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	storeUrl              = "https://online.carrefour.com.tw/tw/"
	baseUrl               = "https://online.carrefour.com.tw"
	defaultNumberOfWorker = 6
	DB_Table              = "carrefour"
)

func main() {
	start := time.Now()

	numWorkers := getNumberOfWorkers()
	log.Printf("Number of workers: %d\n", numWorkers)

	log.Printf("Initiate connection to PostgreSQL")
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Error while initiating connection to database: ", err)
	}

	if err = database.CreateTable(db, DB_Table); err != nil {
		log.Fatal("Error while creating table: ", err)
	}

	doc, err := crawl.GetUrlDocument(storeUrl)
	if err != nil {
		log.Fatalf("Failed to read from store url (%s): %s", storeUrl, err)
	}

	finished := make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	jobPool := concurrent.NewJobPool(numWorkers)

	ctx := gracefulShutdown(context.Background(), func() {
		fmt.Printf("Shutting down gracefully\n")
		wg.Wait()
		close(finished)
	})

	for i := 0; i < numWorkers; i++ {
		jobPool.AddWorker(ctx, wg, db, processPage)
	}

	go func() {
		doc.Find(".top1.left-item").Each(func(i int, selection *goquery.Selection) {
			anchor := selection.Find("a")
			addr, found := anchor.Attr("href")
			if found {
				pageUrl := fmt.Sprintf("%s%s", baseUrl, addr)
				jobPool.Enqueue(pageUrl)
			}
		})
		close(finished)
	}()

	jobPool.Start(ctx)

	<-finished
	elapsed := time.Since(start)
	fmt.Printf("\n\nTook %s to process all categories with %d workers.", elapsed, numWorkers)
}

func getNumberOfWorkers() int {
	numWorkers := flag.Int("numWorkers", defaultNumberOfWorker, fmt.Sprintf("Set the number of workers (default to %d if not specified).", defaultNumberOfWorker))
	flag.Parse()
	if *numWorkers < 1 {
		fmt.Printf("Number of workers has to be at least one. Using default number (%d) instead.\n", defaultNumberOfWorker)
		return defaultNumberOfWorker
	}
	return *numWorkers
}

func gracefulShutdown(c context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(c)
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(ch)

		select {
		case <-ctx.Done():
		case <-ch:
			cancel()
			f()
		}
	}()

	return ctx
}

func processPage(url string, db *sql.DB) error {
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

		productEntry := model.NewProductEntry(name, fmt.Sprintf("%s%s", baseUrl, link), imgLink, price)
		productEntry.PrintProductDetails()
		if err = productEntry.SaveToDB(db, DB_Table); err != nil {
			fmt.Printf("Failed to save %v to db : %s\n", productEntry, err)
		}
	})
	return nil
}
