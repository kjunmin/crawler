package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	visited            map[string]bool
	pages              map[string][]PageNode
	foundRoute         map[string][]PageNode
	baseURL            *url.URL
	targetURL          string
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func configure(rawBaseURL, targetBaseURL string, maxConcurrency int, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}
	targetURL, err := normalizeURL(targetBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't normalize target URL: %v", err)
	}

	return &config{
		visited:            make(map[string]bool),
		pages:              make(map[string][]PageNode),
		foundRoute:         make(map[string][]PageNode),
		baseURL:            baseURL,
		targetURL:          targetURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}
