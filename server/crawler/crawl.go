package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Crawler struct {
	sync.Mutex
	Hostname string
	Run      bool
	Urls     map[string]struct{}
}

// addUrls - Concurrency safe adding urls to crawlers captured urls.
func (c *Crawler) addUrls(urls []string) {
	c.Lock()
	defer c.Unlock()
	for _, url := range urls {
		_, exists := c.Urls[url]
		if exists != false {
			c.Urls[url] = struct{}{}
		}
	}
}

// urlExists - Concurrency safe check if a url has already been captured.
func (c *Crawler) urlExists(url string) bool {
	c.Lock()
	defer c.Unlock()
	_, exists := c.Urls[url]

	return exists
}

// getPageLinks - Retrieves all links on a page.
func (c *Crawler) getPageLinks(body io.Reader) []string {
	var links []string
	tokenizer := html.NewTokenizer(body)

	for {
		tt := tokenizer.Next()

		// Receive stop signal
		if !c.Run {
			return links
		}

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := tokenizer.Token()
			if "a" != token.Data {
				continue
			}

			// Collect a tag hrefs
			for _, attr := range token.Attr {
				if attr.Key != "href" {
					continue
				}
				val := attr.Val

				// Skip internal page links
				if strings.HasPrefix(val, "#") {
					continue
				}

				links = append(links, val)
			}
		}
	}
}

// Crawl - Starts the crawler for a given url.
func (c *Crawler) Crawl(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	c.Run = true
	links := c.getPageLinks(res.Body)
	fmt.Printf("links:: %v", links)
}
