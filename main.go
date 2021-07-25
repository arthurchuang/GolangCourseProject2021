package main

import (
	"database/sql"
	"GoCrawl/concurrent"
	"GoCrawl/crawl"
	//"GoCrawl/model"
	"GoCrawl/db_sql"
	"context"
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
	DB_Table	      = "carrefour"
)

type Database struct{
	DB *sql.DB
}

func timeTrack(start time.Time, name string,numworkers int){
	elapsed := time.Since(start)
	log.Printf("%s took %s with %d worker",name,elapsed,numworkers)
}

func main() {
	numWorkers := getNumberOfWorkers()
	fmt.Printf("Number of workers: %d\n", numWorkers)
	defer timeTrack(time.Now(),"Total Time", numWorkers)

	db:= db_sql.InitDB()
	db_sql.CreateTable(db,DB_Table)
	DD := &Database{DB:db}
	defer db_sql.DeleteAllData(db,DB_Table)
	defer db_sql.QueryDB(db,`SELECT * FROM `+DB_Table)

	doc, err := crawl.GetUrlDocument(storeUrl)
	if err != nil {
		log.Fatalf("Failed to read from store url (%s): %e", storeUrl, err)
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
		jobPool.AddWorker(ctx, wg, DD.processPage)
	}

	jobPool.Start(ctx)

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

	<-finished
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

func (DD *Database) processPage(url string) error {
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

		// productEntry := model.NewProductEntry(name, fmt.Sprintf("%s%s", baseUrl, link), imgLink, price)
		// productEntry.PrintProductDetails()
		// productEntry.SaveToDB()
		db_sql.InsertData(DD.DB,DB_Table,name,link,imgLink,price)
	})
	return nil
}
