package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

const(
	FeedTypeYoutubeChannel = "youtube-channel"
	FeedTypeYoutubePlaylist = "youtube-playlist"
)

type Feed struct {
	ID    uuid.UUID `storm:"unique"`
	Name  string    `req:"name"`
	Type  string    `req:"type"`
	ExtID string    `req:"ext_id"`
	CreatedAt time.Time
}

func (f Feed) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Required),
		validation.Field(&f.Type, validation.In(FeedTypeYoutubeChannel, FeedTypeYoutubePlaylist)),
		validation.Field(&f.ExtID, validation.Required),
	)
}