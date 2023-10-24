package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

var parser = gofeed.NewParser()

func fetchFeed(url string) (gofeed.Feed, error) {

	feed, err := parser.ParseURL(url)
	if err != nil {
		fmt.Printf("Failed to parse feed: %v", err)
	}

	return *feed, nil

}
