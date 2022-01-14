package models

import (
	"github.com/google/uuid"
	"time"
)

type Video struct {
	ID           uuid.UUID `storm:"unique"`
	ExtID        string
	ChannelID    string
	ChannelName  string
	Title        string
	Description  string
	ThumbnailURL string
	UploadedOn   time.Time
	Duration     int
}

type Metadata struct {
	Title        string
	Description  string
	ThumbnailURL string
	UploadTime   time.Time
	Duration     int
}
