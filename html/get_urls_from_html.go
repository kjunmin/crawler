package html

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func constructURL(rawBaseURL, href string) string {
	u, _ := url.Parse(href)
	if u.Host == "" {
		return rawBaseURL + href
	} else {
		return href
	}
}

func GetURLsFromHtml(htmlBody string, baseURL *url.URL) ([]string, error) {
	htmlReader := strings.NewReader(htmlBody)
	nodes, err := html.Parse(htmlReader)

	var urls []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						continue
					}
					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(nodes)

	return urls, err
}
