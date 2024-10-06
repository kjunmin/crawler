package main

import (
	"context"
	"net/url"
	"strings"
)

type PageNode struct {
	nodeType string
	link     string
	label    string
}

func getNodeType(imdbURL string) string {
	if strings.Contains(imdbURL, "title") {
		return "movie"
	}
	return "actor"
}

func (cfg *config) isFirstVisit(currURL string) (isFirstVisit bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, visited := cfg.visited[currURL]; visited {
		return false
	}

	cfg.visited[currURL] = true

	return true
}

func (cfg *config) addPageVisit(cancel context.CancelFunc, label, prevURL, currURL string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.targetURL == currURL {
		cfg.foundRoute[currURL] = append(cfg.pages[prevURL], PageNode{
			label:    label,
			nodeType: getNodeType(currURL),
			link:     currURL,
		})
		cancel()
		return
	}

	if currURL == prevURL {
		cfg.pages[currURL] = []PageNode{{label: label, nodeType: getNodeType(currURL), link: currURL}}
		return
	}

	parsedPrev, _ := url.Parse("//" + prevURL)
	parsedCurr, _ := url.Parse("//" + currURL)
	if !((isActorLink(parsedPrev.Path) && isMovieLink(parsedCurr.Path)) || (isMovieLink(parsedPrev.Path) && isActorLink(parsedCurr.Path))) {
		return
	}

	cfg.pages[currURL] = append(cfg.pages[prevURL], PageNode{
		label:    label,
		nodeType: getNodeType(currURL),
		link:     currURL,
	})

}

func (cfg *config) pagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}
