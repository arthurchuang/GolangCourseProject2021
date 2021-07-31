package model

import (
	"fmt"
)

type ProductEntry struct {
	Name, Link, ImgLink, Price string
}

// NewProductEntry creates and returns a ProductEntry with the passed in product information.
func NewProductEntry(name string, link string, imgLink string, price string) *ProductEntry {
	return &ProductEntry{
		Name:    name,
		Link:    link,
		ImgLink: imgLink,
		Price:   price,
	}
}

func (pe *ProductEntry) String() string {
	return fmt.Sprintf("name: %s\nlink: %s\nimage: %s\nprice: %s\n\n", pe.Name, pe.Link, pe.ImgLink, pe.Price)
}
