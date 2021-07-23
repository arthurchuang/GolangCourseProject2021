package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"
	"GoCrawl/db_sql"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/PuerkitoBio/goquery"
)

const (
	storeUrl = "https://online.carrefour.com.tw/tw/"
	baseUrl  = "https://online.carrefour.com.tw"
	USER = "user"
	PASSWORD= "123456"
	DATABASE = "postgres"
	HOST = "localhost"
	port = 5432
	DB_NAME="carrefour"
)

func timeTrack(start time.Time, name string){
	elapsed := time.Since(start)
	log.Printf("%s took %s",name,elapsed)
}

type Database struct{
	DB *sql.DB
}
func main() {
	defer timeTrack(time.Now(),"Total Time")
	db, err := sql.Open(
		"postgres",fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",HOST, USER,PASSWORD,DATABASE))
	if err = db.Ping(); err!=nil{
		log.Println("Fail to connect to db")
		panic(err)
	}
	fmt.Println("Successfully created connection to database")
	
	env := &Database{DB:db}
	db_sql.QueryDB(env.DB,`SELECT * FROM `+DB_NAME)
	db_sql.CreateTable(env.DB,DB_NAME)
	defer db_sql.QueryDB(env.DB,`SELECT * FROM `+DB_NAME)
	defer db_sql.DeleteAllData(env.DB,DB_NAME)
	//db_sql.InsertData(env.DB, DB_NAME,"harry","link","imgLink","1500")

	res, err := http.Get(storeUrl)
	if err != nil {
		log.Fatalf("Error getting store website: %e", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}
	worker_number := 10
	var wg sync.WaitGroup
	job_list := make(chan string,1)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	go func(){
		doc.Find(".top1.left-item").Each(func(i int, selection *goquery.Selection) {
		anchor := selection.Find("a")
		addr, found := anchor.Attr("href")
		if found {
			pageUrl := fmt.Sprintf("%s%s", baseUrl, addr)
			job_list<- pageUrl
 			}
		})
		defer close(job_list)
	}()
	defer wg.Wait()
	wg.Add(1)
	for w:=1; w<=worker_number;w++{
		go Worker(w,job_list,&wg,env.DB)
	}
}

func Worker(id int,jobs <-chan string, wg *sync.WaitGroup, db *sql.DB){
	defer wg.Done()
	for job:= range jobs{
		fmt.Println("worker number:",id)
		if err := processPage(job,db); err != nil {
		fmt.Printf("Error while processing page (%s) : %e", job, err)
			}
	}
}

func processPage(url string,db *sql.DB) error {
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
		//db_sql.InsertData(env.db, DB_NAME, name, link, imgLink, price)
		saveEntry(db,name, link, imgLink, price)
	})
	return nil
}

func saveEntry(db *sql.DB,name string, link string, imgLink string, price string) {
	// fmt.Printf("product name: %s\n", name)
	// fmt.Printf("product link: %s%s\n", baseUrl, link)
	// fmt.Printf("product image: %s\n", imgLink)
	// fmt.Printf("product price: %s", price)
	// fmt.Printf("\n\n")
	db_sql.InsertData(db,DB_NAME,name,link,imgLink,price)

	// TODO : save to DB
}
