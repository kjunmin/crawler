package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func main() {

	inputArgs := os.Args[1:]

	if len(inputArgs) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(inputArgs) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := inputArgs[0]

	fmt.Printf("starting crawl of: %v", baseUrl)

	pages := make(map[string]int)
	crawlPage(baseUrl, baseUrl, pages)

	PrettyPrint(pages)
}
