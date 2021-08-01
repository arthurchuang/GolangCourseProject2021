package main

import (
	"GoCrawl/concurrent"
	"GoCrawl/crawl"
	"GoCrawl/database"
	"GoCrawl/model"
	"context"
	"database/sql"
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
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
	dbHost     = "localhost"

	storeUrl    = "https://online.carrefour.com.tw/tw/"
	baseUrl     = "https://online.carrefour.com.tw"
	dbTableName = "carrefour"
)

func main() {
	start := time.Now()

	numWorkers := concurrent.GetNumberOfWorkers()
	log.Printf("Number of workers: %d\n", numWorkers)

	log.Println("Initiating connection to PostgreSQL database")
	db, err := database.InitDB(dbHost, dbName, dbUser, dbPassword)
	if err != nil {
		log.Fatalf("Error while initiating connection to database: %s", err)
	}

	log.Println("Preparing table for recording product entries")
	if err = database.CreateTableIfNotExist(db, dbTableName); err != nil {
		log.Fatalf("Error while creating table: %s", err)
	}

	if err = database.TruncateTable(db, dbTableName); err != nil {
		log.Fatalf("Error while deleting previous content in %s table : %s", dbTableName, err)
	}
	log.Println("Database setup completed")

	doc, err := crawl.GetUrlDocument(storeUrl)
	if err != nil {
		log.Fatalf("Failed to read from store url (%s): %s", storeUrl, err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	jobPool := concurrent.NewJobPool(numWorkers)

	ctx := gracefulShutdown(context.Background(), func() {
		fmt.Printf("Shutting down gracefully\n")
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
		jobPool.NoMoreInput()
	}()

	jobPool.Start(ctx)

	wg.Wait()
	elapsed := time.Since(start)

	numSaved, err := database.GetElementCounts(db, dbTableName)
	if err != nil {
		log.Fatalf("Failed to get number of product entries saved: %s", err)
	}
	fmt.Printf("Saved %d product entries with %d workers in %s\n", numSaved, numWorkers, elapsed)
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
		fmt.Print(productEntry)
		if err = database.SaveProductEntry(db, dbTableName, productEntry); err != nil {
			fmt.Printf("Failed to save %v to db : %s\n", productEntry, err)
		}
	})
	return nil
}
