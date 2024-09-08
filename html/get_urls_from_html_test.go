package html

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		htmlBody string
		expected []string
	}{
		{
			name:     "HTML parse",
			inputURL: "https://blog.boot.dev",
			htmlBody: `
	var urls []string

				<html>
					<body>
						<a href="/path/one">
							<span>Boot.dev</span>
						</>
						<a href="https://other.com/path/one">
							<span>Boot.dev</span>
						</>
					</body>
				</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseUrl, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: couldn't parse url %v", i, tc.name, tc.inputURL)
			}
			actual, err := GetURLsFromHtml(tc.htmlBody, baseUrl)
			if err != nil {
				t.Errorf("Test %v - %s FAIL actual %v, expected %v", i, tc.name, actual, tc.expected)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL actual %v, expected %v", i, tc.name, actual, tc.expected)
			}
		})
	}
}
