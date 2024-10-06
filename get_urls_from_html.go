package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func isActorLink(path string) bool {
	// e.g. Actor link on IMDB: www.imdb.com/name/nm9423355
	parsed, _ := url.Parse(path)
	re := regexp.MustCompile(`^/name/nm\d{7}/?$`)
	match := re.MatchString(parsed.Path)
	if match {
		return true
	}
	return false
}

func isMovieLink(path string) bool {
	// e.g. Title link on IMDB: www.imdb.com/title/tt19715796
	parsed, _ := url.Parse(path)
	re := regexp.MustCompile(`^/title/tt\d{8}/?$`)
	match := re.MatchString(parsed.Path)
	if match {
		return true
	}
	return false
}

func getActorOrMovieName(node *html.Node, name *string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "itemprop" && attr.Val == "url" {
				*name = node.FirstChild.Data
				return
			}
		}
	}
	if node.Type == html.ElementNode && node.Data == "span" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "hero__primary-text" {
				*name = node.FirstChild.Data
				return
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		getActorOrMovieName(child, name)
	}
}

func getMovieCastNode(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "table" {
		for _, attributes := range node.Attr {
			if attributes.Key == "class" && attributes.Val == "cast_list" {
				return node
			}
		}
	}
	return nil
}

func getContainingNode(node *html.Node, resNode **html.Node, fn func(*html.Node) *html.Node) {
	res := fn(node)
	if res != nil {
		*resNode = res
		return
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		getContainingNode(child, resNode, fn)
	}
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) (string, []string, error) {
	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't parse HTML: %v", err)
	}

	var name string
	var urls []string
	var resNode *html.Node
	getActorOrMovieName(doc, &name)
	getContainingNode(doc, &resNode, getMovieCastNode)
	if resNode == nil {
		resNode = doc
	}

	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, anchor := range node.Attr {
				if anchor.Key == "href" {
					href, err := url.Parse(anchor.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", anchor.Val, err)
						continue
					}

					if !(isActorLink(href.Path) || isMovieLink(href.Path)) {
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseNodes(child)
		}
	}
	traverseNodes(resNode)

	return name, urls, nil
}
