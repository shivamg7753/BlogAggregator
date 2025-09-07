package jobs

import (
	"blogAggregator/internal/database"
	"blogAggregator/internal/models"
	"blogAggregator/internal/rss"
	"fmt"
	"time"
)

func StartFeedUpdater(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		fmt.Println("running feed updater.....")
		var feeds []models.Feed
		if err:=database.DB.Find(&feeds).Error;err!=nil{
			fmt.Println("could not fetch feeds:" ,err)
			continue
		}
		for _,feed:=range feeds{
			fmt.Printf("refreshing feed: %s\n",feed.Title)
			if err:=rss.FetchAndStoreFeed(feed);err!=nil{
				fmt.Println("error refreshing feed: ",err)
			}
		}
	}
}
