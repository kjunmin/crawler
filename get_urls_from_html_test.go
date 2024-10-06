package main

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestIsActorLink(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected bool
	}{
		{
			name:     "Incorrect Actor link",
			inputURL: "http://www.imdb.com/title/tt19715796",
			expected: false,
		},
		{
			name:     "Correct actor link",
			inputURL: "http://www.imdb.com/name/nm9423355/",
			expected: true,
		},
		{
			name:     "Correct actor link",
			inputURL: "http://www.imdb.com/name/nm9423355",
			expected: true,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("TEST %v FAIL: Unable to parse input URL %v", i, tc.inputURL)
			}
			fmt.Println(parsed.Path)
			isActor := isActorLink(parsed.Path)
			if isActor != tc.expected {
				t.Errorf("TEST %v FAIL: Incorrect link determination for %v. Expected %v, got %v", i, tc.inputURL, tc.expected, isActor)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="https://blog.boot.dev">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
