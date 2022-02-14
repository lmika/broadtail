package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
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
	VideoID   string
	Title     string
	Link      string
	Published time.Time
}
