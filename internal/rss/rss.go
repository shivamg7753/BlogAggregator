package rss

import (
	"blogAggregator/internal/database"
	"blogAggregator/internal/models"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

func FetchAndStoreFeed(feed models.Feed) error {
	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feed.URL)
	if err != nil {
		return fmt.Errorf("failed to parse feed %s:%w", feed.URL, err)
	}
	now := time.Now().UTC()
  database.DB.Model(&feed).Update("last_fetched", &now)
	for _, item := range parsedFeed.Items {
		published := time.Now().UTC()
    if item.PublishedParsed!=nil{
			published=*item.PublishedParsed
		}
		post := models.Post{
			Title: item.Title,
			Link: item.Link,
			Content: item.Content,
			Published: published,
			FeedId: feed.ID,
		}

		result:=database.DB.Where("link=?",post.Link).FirstOrCreate(&post)
		if result.Error!=nil{
			return result.Error
		}
	}
	return nil
}
