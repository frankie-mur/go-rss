package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/google/uuid"
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

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
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
		//TODO: Description and published at not being inserted correctly
		feed_time := sql.NullTime{}
		if f.PubDate == "" {
			feed_time.Time = time.Now()
			feed_time.Valid = false
		}
		time_parsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", f.PubDate)
		if err != nil {
			log.Printf("Failed to parse time: %v", err)

		}
		feed_time.Time = time_parsed
		feed_time.Valid = true
		post, err := db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       f.Title,
				Url:         f.Link,
				Description: f.Description,
				PublishedAt: feed_time,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			//Want to ignore duplicate key errors as they are Posts already in db
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Failed with error: %v", err)
		}

		log.Printf("Created Post with id %v, and title %s", post.ID, post.Title)
	}

	log.Printf("Fetched feed with title %s, with %d posts", fetched_feed.Channel.Title, len(fetched_feed.Channel.Item))
}
