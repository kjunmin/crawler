package url

import (
	"net/url"
	"strings"
)

func trimLastChar(s, char string) string {
	if strings.HasSuffix(s, char) {
		return s[:len(s)-len(char)]
	}
	return s
}

func NormalizeURL(inputURL string) (normalizedURL string, err error) {

	u, err := url.Parse(inputURL)

	normalizedURL = trimLastChar(u.Host+u.Path, "/")
	return normalizedURL, err
}
