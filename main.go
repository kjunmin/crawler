package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 5 {
		fmt.Println("not enough arguments provided")
		fmt.Println("usage: crawler <baseURL> <maxConcurrency> <maxPages>")
		return
	}
	if len(os.Args) > 5 {
		fmt.Println("too many arguments provided")
		return
	}
	rawBaseURL := os.Args[1]
	targetBaseURL := os.Args[2]
	maxConcurrencyString := os.Args[3]
	maxPagesString := os.Args[4]

	maxConcurrency, err := strconv.Atoi(maxConcurrencyString)

	if err != nil {
		fmt.Printf("Error - maxConcurrency: %v", err)
		return
	}
	maxPages, err := strconv.Atoi(maxPagesString)
	if err != nil {
		fmt.Printf("Error - maxPages: %v", err)
		return
	}

	cfg, err := configure(rawBaseURL, targetBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure: %v", err)
		return
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	ctx, cancel := context.WithCancel(context.Background())
	cfg.wg.Add(1)
	go cfg.crawlPage(ctx, cancel, rawBaseURL, rawBaseURL)
	cfg.wg.Wait()

	printReport(cfg.pages, cfg.foundRoute, rawBaseURL)
}
