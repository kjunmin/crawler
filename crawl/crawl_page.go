package crawl

import (
	"fmt"
	"net/url"

	crawlhtml "github.com/kjunmin/crawler/html"
	crawlurl "github.com/kjunmin/crawler/url"
)

func CrawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {

	u, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Errorf("Error parsing URL %v", rawBaseURL)
		return
	}

	baseHost := u.Host

	u2, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Error parsing URL %v", rawCurrentURL)
		return
	}

	currentHost := u2.Host

	if currentHost != baseHost {
		fmt.Println("Base host differs... returning")
		return
	}

	normalizedRawCurrentURL, err := crawlurl.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Unable to normalize URL %v", rawCurrentURL)
		return
	}

	if _, exists := pages[normalizedRawCurrentURL]; exists {
		fmt.Printf("%v already crawled %v times", normalizedRawCurrentURL, pages[normalizedRawCurrentURL])
		return
	}

	htmlBody, err := crawlhtml.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Errorf("Unable to get HTML for %v", rawCurrentURL)
		return
	}

	urls, err := crawlhtml.GetURLsFromHtml(htmlBody, normalizedRawCurrentURL)
	if err != nil {
		fmt.Errorf("Unable to get URLS from html body for URL %v", normalizedRawCurrentURL)
	}

	for _, url := range urls {
		normalizedUrl, err := crawlurl.NormalizeURL(url)
		if err != nil {
			continue
		}
		fmt.Println(normalizedUrl)
		if _, exists := pages[normalizedUrl]; exists {
			fmt.Printf("URL %v already exists in page table", normalizedUrl)
		} else {
			fmt.Printf("Crawling url %v...", normalizedUrl)
			pages[normalizedUrl]++
			CrawlPage(rawBaseURL, normalizedUrl, pages)
		}
	}

}
