package model

import "fmt"

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
	fmt.Printf("product name: %s\n", pe.name)
	fmt.Printf("product link: %s\n", pe.link)
	fmt.Printf("product image: %s\n", pe.imgLink)
	fmt.Printf("product price: %s\n\n", pe.price)
}

// SaveToDB saves the productEntry to the database.
func (pe *productEntry) SaveToDB() {
	// TODO
}
