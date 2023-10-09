package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
)

func initScraper(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scrapping %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("Error getting feeds: %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(&wg, db, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Failed to mark feed as fetched: %v", err)
		return
	}

	fetched_feed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Failed to fetch feed: %v", err)
		return
	}
	for _, f := range fetched_feed.Channel.Item {
		log.Printf("Found post %s", f.Title)
	}

	log.Printf("Fetched feed with title %s, with %d posts", fetched_feed.Channel.Title, len(fetched_feed.Channel.Item))
}
