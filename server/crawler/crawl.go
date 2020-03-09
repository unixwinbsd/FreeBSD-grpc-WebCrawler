package crawler

import (
	"github.com/gocolly/colly/v2"
	"log"
	"net/url"
	"strings"
	"sync"
)

type Crawler struct {
	sync.Mutex
	Hostname string
	Run      bool
	Urls     map[string]map[string]struct{}
}

// addUrls - Concurrency safe adding urls to crawlers captured urls.
func (crawler *Crawler) addUrl(hostname string, url string) {
	crawler.Lock()
	defer crawler.Unlock()

	_, exists := crawler.Urls[hostname]
	if !exists {
		crawler.Urls[hostname] = map[string]struct{}{}
	}

	crawler.Urls[hostname][url] = struct{}{}
}

// formatUrl - Remove trailing slash, fix partials and handle links external domain.
func (crawler *Crawler) formatUrl(linkUrl string) string {
	u := strings.TrimRight(linkUrl, "/")

	// Trim off http/s from url to normalize
	if strings.HasPrefix(linkUrl, "http") {
		u = strings.Split(linkUrl, "://")[1]
	}

	// Fix internal links with full patch
	if strings.HasPrefix(u, "/") {
		u = crawler.Hostname + u
	}

	// No external links
	if !strings.Contains(u, crawler.Hostname) {
		return ""
	}

	// Page/content links
	if strings.HasPrefix(u, "#") {
		return ""
	}

	return u
}

// Crawl - Starts the crawler for a given url.
func (crawler *Crawler) Crawl(crawlUrl string) {
	u, err := url.Parse(crawlUrl)
	if err != nil {
		log.Fatal(err)
	}

	hostname := u.Hostname()
	crawler.Hostname = hostname

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains(hostname),
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !crawler.Run {
			log.Println("crawler stopped, exiting request")
			return
		}

		link := e.Attr("href")
		formatted := crawler.formatUrl(link)
		if formatted != "" {
			crawler.addUrl(hostname, formatted)
		}

		// Visit link found on page on a new thread
		e.Request.Visit(link)
	})

	// Start scraping on url
	c.Visit(crawlUrl)

	// Wait until threads are finished
	c.Wait()
	log.Println(crawlUrl, "finished crawl")
}
