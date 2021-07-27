package model

import (
	"GoCrawl/database"
	"database/sql"
	"log"
)

type productEntry struct {
	name, link, imgLink, price string
}

// NewProductEntry creates and returns a productEntry with the passed in product information.
func NewProductEntry(name string, link string, imgLink string, price string) *productEntry {
	return &productEntry{
		name:    name,
		link:    link,
		imgLink: imgLink,
		price:   price,
	}
}

// PrintProductDetails prints out the details of the productEntry.
func (pe *productEntry) PrintProductDetails() {
	log.Printf("product name: %s\n", pe.name)
	log.Printf("product link: %s\n", pe.link)
	log.Printf("product image: %s\n", pe.imgLink)
	log.Printf("product price: %s\n\n", pe.price)
}

// SaveToDB saves the productEntry to the database.
func (pe *productEntry) SaveToDB(db *sql.DB, table string) error {
	return database.InsertData(db, table, pe.name, pe.link, pe.imgLink, pe.price)
}
