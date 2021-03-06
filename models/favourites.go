package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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

	// Dealing with old syle video references
	if len(refId) == 1 {
		return VideoRef{Source: YoutubeVideoRefSource, ID: refId[0]}, nil
	}

	if len(refId) != 2 {
		return VideoRef{}, errors.Errorf("invalid manual ref: ")
	}

	return VideoRef{Source: VideoRefSource(refId[0]), ID: refId[1]}, nil
}

func (vr VideoRef) String() string {
	return fmt.Sprintf("%v:%v", vr.Source, vr.ID)
}

type FavouriteOrigin struct {
	Type OriginType `json:"type"`
	ID   string     `json:"id"`
}

type VideoRefSource string

const (
	YoutubeVideoRefSource  = "youtube"
	AppleDevVideoRefSource = "apple-dev"
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
