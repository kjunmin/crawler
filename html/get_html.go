package html

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetHTML(rawURL string) (string, error) {

	resp, err := http.Get(rawURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("status code %v", resp.StatusCode)
	}

	if !strings.Contains(resp.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("invalid content type %v", resp.Header.Get("content-type"))
	}

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contentBytes), nil
}
