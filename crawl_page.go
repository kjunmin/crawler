package main

import (
	"fmt"
	"net/url"
	"sync"

	crawlhtml "github.com/kjunmin/crawler/html"
	crawlurl "github.com/kjunmin/crawler/url"
)

type config struct {
	pages              map[string]int
	maxPages           int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func configure(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {

	u, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing URL %v", rawBaseURL)
	}

	return &config{
		pages:              make(map[string]int),
		maxPages:           maxPages,
		baseURL:            u,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}, nil
}

func (cfg *config) addPageVisit(normalizedUrl string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.pages[normalizedUrl]; visited {
		cfg.pages[normalizedUrl]++
		return false
	}

	cfg.pages[normalizedUrl] = 1
	return true
}

func (cfg *config) checkMaxPagesReached() (isReached bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return true
	}
	return false
}

func (cfg *config) CrawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	isReached := cfg.checkMaxPagesReached()
	if isReached {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Error parsing URL %v", rawCurrentURL)
		return
	}

	if currentURL.Hostname() != cfg.baseURL.Hostname() {
		fmt.Println("Base host differs... returning")
		return
	}

	normalizedRawCurrentURL, err := crawlurl.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Unable to normalize URL %v", rawCurrentURL)
		return
	}

	isFirst := cfg.addPageVisit(normalizedRawCurrentURL)
	if !isFirst {
		return
	}

	fmt.Printf("Crawling %s\n", rawCurrentURL)

	htmlBody, err := crawlhtml.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Unable to get HTML for %v", rawCurrentURL)
		return
	}

	nextURLs, err := crawlhtml.GetURLsFromHtml(htmlBody, cfg.baseURL)
	if err != nil {
		fmt.Errorf("Unable to get URLS from html body for URL %v", normalizedRawCurrentURL)
	}

	for _, nextURL := range nextURLs {
		cfg.wg.Add(1)
		go cfg.CrawlPage(nextURL)
	}

}
