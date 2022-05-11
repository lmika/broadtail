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
	ID              uuid.UUID `storm:"id"`
	Name            string    `req:"name"`
	Type            string    `req:"type"`
	ExtID           string    `req:"ext_id"`
	TargetDir       string    `req:"target_dir"`
	Ordering        string    `req:"ordering"`
	CheckForUpdates bool      `req:"check_for_updates"`
	CreatedAt       time.Time
	LastUpdatedAt   time.Time
}

func (f Feed) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Type, validation.In(FeedTypeYoutubeChannel, FeedTypeYoutubePlaylist, FeedTypeAppleDev)),
		validation.Field(&f.ExtID, validation.Required),
		validation.Field(&f.Ordering, validation.In(ChronologicalFeedItemOrdering, AlphabeticalFeedItemOrdering)),
	)
}

const (
	ChronologicalFeedItemOrdering = "pub-desc"
	AlphabeticalFeedItemOrdering  = "title-asc"
)

type FeedHints struct {
	Ordering        string
	CheckForUpdates bool
}
