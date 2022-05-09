package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

const (
	FeedTypeYoutubeChannel  = "youtube-channel"
	FeedTypeYoutubePlaylist = "youtube-playlist"
)

type Feed struct {
	ID            uuid.UUID `storm:"id"`
	Name          string    `req:"name"`
	Type          string    `req:"type"`
	ExtID         string    `req:"ext_id"`
	TargetDir     string    `req:"target_dir"`
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

func (f Feed) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Type, validation.In(FeedTypeYoutubeChannel, FeedTypeYoutubePlaylist)),
		validation.Field(&f.ExtID, validation.Required),
	)
}

type FeedItem struct {
	ID        uuid.UUID `storm:"id"`
	EntryID   string    `storm:"unique"`
	FeedID    uuid.UUID `storm:"index"`
	Title     string
	Link      string
	Published time.Time
}

func (fi FeedItem) VideoRef() VideoRef {
	p, err := ParseVideoRef(fi.EntryID)
	if err == nil {
		return p
	}

	// Likely to be an old entry ID
	return VideoRef{
		Source: YoutubeVideoRefSource,
		ID:     fi.EntryID,
	}
}

func (fi *FeedItem) SetVideoRef(vr VideoRef) {
	fi.EntryID = vr.String()
}

type RecentFeedItem struct {
	FeedItem    FeedItem
	Feed        Feed
	FavouriteID string
}
