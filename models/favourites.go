package models

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strings"
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

func ParseVideoRef(str string) (VideoRef, error) {
	refId := strings.SplitN(str, ":", 2)
	if len(refId) != 2 {
		return VideoRef{}, errors.Errorf("invalid manual ref: ")
	}

	// TODO:
	if refId[0] != YoutubeVideoRefSource {
		return VideoRef{}, errors.Errorf("unrecognised source")
	}

	return VideoRef{Source: YoutubeVideoRefSource, ID: refId[1]}, nil
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
	ManualOriginType   = "manual"
	FeedItemOriginType = "feed-item"
)

type FavouriteWithOrigin struct {
	Favourite   Favourite
	OriginTitle string
	OriginURL   string
}
