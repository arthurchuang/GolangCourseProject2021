package crawl

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// GetUrlDocument returns a *goquery.Document containing the content of the website with the given url.
func GetUrlDocument(url string) (*goquery.Document, error) {
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

	return doc, nil
}
