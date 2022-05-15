package models

import (
	"github.com/google/uuid"
	"time"
)

type FeedItem struct {
	ID        uuid.UUID `storm:"id"`
	VideoRef  VideoRef  `storm:"unique"`
	FeedID    uuid.UUID `storm:"index"`
	Title     string
	Link      string
	Published time.Time
}

type RecentFeedItem struct {
	FeedItem    FeedItem
	Feed        Feed
	FavouriteID string
	Downloaded  bool
}

// FetchedFeedItem is a feed item fetched from a RSS source
type FetchedFeedItem struct {
	VideoRef    VideoRef
	Title       string
	Description string
	Link        string
	Published   time.Time
}
