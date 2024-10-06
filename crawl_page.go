package main

import (
	"context"
	"fmt"
	"net/url"
)

func constructCrawlURL(rawBaseURL string) string {
	normalized, _ := normalizeURL(rawBaseURL)
	if isMovieLink(normalized) {
		joinedURL, _ := url.JoinPath(normalized, "/fullcredits")
		return joinedURL
	}
	return normalized
}

func (cfg *config) crawlPage(ctx context.Context, cancel context.CancelFunc, rawPrevURL, rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Target found, worker stopping...")
			return
		default:
			if cfg.pagesLen() >= cfg.maxPages {
				return
			}

			currentURL, err := url.Parse(rawCurrentURL)
			if err != nil {
				fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
				return
			}

			// skip other websites
			if currentURL.Hostname() != cfg.baseURL.Hostname() {
				return
			}

			normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
			if err != nil {
				fmt.Printf("Error - normalizedURL: %v", err)
				return
			}

			isFirst := cfg.isFirstVisit(normalizedCurrentURL)
			if !isFirst {
				return
			}

			normalizedPreviousURL, err := normalizeURL(rawPrevURL)
			if err != nil {
				fmt.Printf("Error - normalizedURL: %v", err)
				return
			}

			fmt.Printf("crawling %s", rawCurrentURL)

			htmlBody, err := getHTML(constructCrawlURL(rawCurrentURL))
			if err != nil {
				fmt.Printf("Error - getHTML: %v", err)
				return
			}

			label, nextURLs, err := getURLsFromHTML(htmlBody, cfg.baseURL)
			if err != nil {
				fmt.Printf("Error - getURLsFromHTML: %v", err)
				return
			}
			fmt.Printf(" Name %s \n", label)

			cfg.addPageVisit(cancel, label, normalizedPreviousURL, normalizedCurrentURL)

			for _, nextURL := range nextURLs {
				cfg.wg.Add(1)
				go cfg.crawlPage(ctx, cancel, normalizedCurrentURL, nextURL)
			}
		}
	}
}
