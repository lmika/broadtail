package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

const (
	FeedTypeYoutubeChannel  = "youtube-channel"
	FeedTypeYoutubePlaylist = "youtube-playlist"
	FeedTypeAppleDev        = "apple-dev"
)

type Feed struct {
	ID            uuid.UUID        `storm:"id"`
	Name          string           `req:"name"`
	Type          string           `req:"type"`
	ExtID         string           `req:"ext_id"`
	TargetDir     string           `req:"target_dir"`
	Ordering      FeedItemOrdering `req:"ordering"`
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

func (f Feed) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Type, validation.In(FeedTypeYoutubeChannel, FeedTypeYoutubePlaylist, FeedTypeAppleDev)),
		validation.Field(&f.ExtID, validation.Required),
		validation.Field(&f.Ordering, validation.In(ChronologicalFeedItemOrdering, AlphabeticalFeedItemOrdering)),
	)
}

type FeedItemOrdering int

const (
	ChronologicalFeedItemOrdering FeedItemOrdering = 0
	AlphabeticalFeedItemOrdering                   = 1
)

type FeedHints struct {
	Ordering FeedItemOrdering
}
