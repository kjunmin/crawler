package main

import (
	"fmt"
	"os"
	"strconv"
)

func PrettyPrint(v map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %v\n", baseURL)
	fmt.Println("=============================")

	for item, count := range v {
		fmt.Printf("Found %v internal links to %v\n", count, item)
	}

}

func main() {

	inputArgs := os.Args[1:]

	if len(inputArgs) < 3 {
		fmt.Println("Please supply 3 arguments [website] [max_concurrency] [max_pages]")
		os.Exit(1)
	}
	if len(inputArgs) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(inputArgs[1])
	if err != nil {
		fmt.Errorf("Invalid input for concurrency (should be a number)")
		return
	}
	maxPages, err := strconv.Atoi(inputArgs[2])
	if err != nil {
		fmt.Errorf("Invalid input for maxPages (should be a number)")
		return
	}

	baseUrl := inputArgs[0]

	fmt.Printf("starting crawl of: %v", baseUrl)

	cfg, err := configure(baseUrl, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - configure %v", err)
		return
	}

	cfg.wg.Add(1)
	go cfg.CrawlPage(baseUrl)
	cfg.wg.Wait()

	PrettyPrint(cfg.pages, baseUrl)
}
