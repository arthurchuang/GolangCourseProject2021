package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func main() {
	baseURL := "https://online.carrefour.com.tw/zh/homepage"
	response, err := http.Get(baseURL)
	if err != nil {
		log.Fatal("Unable to parse from the baseURL: ", err)
	}
	body, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Fatal("Failed to read from HTML's body: ", errRead)
	}
	defer response.Body.Close()

	// Install the html package, using `go get golang.org/x/net/html`
	reader := strings.NewReader(string(body))
	tokenizer := html.NewTokenizer(reader)

	// Scan user's input for category name
	fmt.Println("Enter the category name (i.e.生鮮食品/冷凍食品/飲料零食): ")
	var first string
	fmt.Scanln(&first)

	for {

		// Iterate through each token and check the token type to find anchor tags (link).
		tt := tokenizer.Next()
		t := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			log.Fatal(err)
		case html.StartTagToken:
			// Check if the token is an <a> tag.
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			// Extract the href value, if there is one
			ok, encodedValue := getHref(t)
			if !ok {
				continue
			}
			// Decode the URL strings found
			decodedValue, err := url.QueryUnescape(encodedValue)
			if err != nil {
				log.Fatal(err)
				return
			}
			// Print out all relevant URLs for the category specified by user input
			if strings.Contains(decodedValue, first) {
				fmt.Printf("URL: %q\n", decodedValue)
			}
		}
	}
}
