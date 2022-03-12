package models

import (
	"github.com/google/uuid"
	"time"
)

type Favourite struct {
	ID         uuid.UUID       `json:"id" storm:"id"`
	VideoRef   VideoRef        `json:"videoRef" storm:"unique"`
	Origin     FavouriteOrigin `json:"origin" storm:"index"`
	Title      string          `json:"title"`
	Link       string          `json:"link"`
	Published  time.Time       `json:"published"`
	Favourited time.Time       `json:"favourited"`
}

type VideoRef struct {
	Source VideoRefSource `json:"source"`
	ID     string         `json:"id"`
}

type FavouriteOrigin struct {
	Type OriginType `json:"type"`
	ID   string     `json:"id"`
}

type VideoRefSource string

const (
	YoutubeVideoRefSource = "youtube"
)

type OriginType string

const (
	FeedItemOriginType = "feed-item"
)
