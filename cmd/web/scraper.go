package main

import (
	"context"
	"database/sql"
	"fmt"
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

	//Check if it's been marked as fetched in the last hour if so return immediately
	if !feed.LastFetchedAt.Time.Add(time.Hour).Before(time.Now()) {
		fmt.Printf("Feed %v has been marked as fetched in the last hour\n", feed.ID)
		return
	}

	fetchedFeed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Failed to fetch feed: %v", err)
		return
	}

	_, err = db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Failed to mark feed as fetched: %v", err)
		return
	}

	for _, f := range fetchedFeed.Items {
		post, err := db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       f.Title,
				Url:         f.Link,
				Description: f.Description,
				PublishedAt: sql.NullTime{
					Time:  *f.PublishedParsed,
					Valid: true,
				},
				FeedID: feed.ID,
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

	log.Printf("Fetched feed with title %s, with %d posts", fetchedFeed.Title, len(fetchedFeed.Items))
}
