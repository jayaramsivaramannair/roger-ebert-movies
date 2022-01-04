package main


import (
	"fmt"
	"strings"
	"net/http"
	"golang.org/x/net/html"
)

func getSource(t html.Token)(ok bool, url string) {
	for _, a := range t.Attr {
		if a.Key == "src" {
			url = a.Val
			ok = true
		}
	}

	return ok, url
}


func crawl(baseURL string, ch chan string) {
	resp, err := http.Get(baseURL)

	if err != nil {
		fmt.Println("Error: Failed to crawl:", baseURL)
	}

	b := resp.Body

	fmt.Println(b)

	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
			case tt == html.ErrorToken:
				return 

			case tt == html.SelfClosingTagToken:
				t := z.Token()

			isImage := t.Data == "img"
			if !isImage {
				continue
			}

			ok, src := getSource(t)

			if !ok {
				continue
			}

			hasProto := strings.Index(src, "http") == 0

			if hasProto{
				ch <- src
			}
		}
	}
}

func main() {


	baseURL := "https://www.rogerebert.com/great-movies"
	imageURLs := make(chan string)

	

	go crawl(baseURL, imageURLs)

	for src := range imageURLs {
		fmt.Println(src)
	}

}