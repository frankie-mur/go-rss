package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetchFeeds(t *testing.T) {
	// Mock RSS feed content
	mockRSS := `
		<?xml version="1.0"?>
		<rss version="2.0">
			<channel>
				<title>Mock RSS Feed</title>
			</channel>
		</rss>
	`
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockRSS)
	}))
	defer server.Close()
	t.Run("fetch returns a Rss", func(t *testing.T) {
		var url = server.URL
		got, err := fetchFeed(url)
		if err != nil {
			t.Errorf("Function called with error: %v", err)
		}
		want := gofeed.Feed{}
		if reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	// t.Run("should return error for unvalid URL", func(t *testing.T) {
	// 	var url = server.URL
	// 	got, err := fetchFeed(url + "/error")
	// 	if err != nil {
	// 		t.Errorf("Function called with error: %v", err)
	// 	}

	// })
}
