package main

import (
	"fmt"
)

type Page struct {
	URL          string
	pageNodeList []PageNode
}

func printReport(pages map[string][]PageNode, foundRoute map[string][]PageNode, baseURL string) {
	fmt.Printf(`
=============================
  REPORT for %s
=============================
`, baseURL)

	sortedPages := sortPages(pages)
	for _, page := range sortedPages {
		url := page.URL
		pageNodeList := page.pageNodeList
		fmt.Printf("Shortest path from a to %s is %d links \n", url, len(pageNodeList))
		for _, pageNode := range pageNodeList {
			fmt.Printf("-> %v . %v - %v\n", pageNode.nodeType, pageNode.label, pageNode.link)
		}
		fmt.Println()
	}
}

func sortPages(pages map[string][]PageNode) []Page {
	pagesSlice := []Page{}
	for url, pageNodeList := range pages {
		pagesSlice = append(pagesSlice, Page{URL: url, pageNodeList: pageNodeList})
	}
	return pagesSlice
}
