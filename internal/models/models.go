package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}


type Feed struct{
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	URL string `gorm:"uniqueIndex;not null" json:"url"`
  CreatedAt time.Time `json:"created_at"`
	LastFetched *time.Time `json:"last_fetched"`
}

type Post struct{
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Link string `gorm:"uniqueIndex;not null" json:"link"`
	Content string `json:"content"`
	Published time.Time `json:"published"`
	FeedId uint `json:"feed_id"`
}


type Subscription struct{
	ID uint `gorm:"primaryKey" json:"id"`
	UserID uint `json:"user_id"`
	FeedID uint `json:"feed_id"`
}