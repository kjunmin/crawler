package main

import (
	"net/url"
	"strings"
)

func truncateIfContainsSlash(s string) string {
	if idx := strings.LastIndex(s, "/"); idx != -1 {
		return s[:idx]
	}
	return s
}

func normalizeURL(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)

	trimmedPath := strings.TrimSuffix(parsedURL.Scheme+"://"+parsedURL.Host+parsedURL.Path, "/")

	return trimmedPath, err
}
