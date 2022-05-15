package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	VideoRef     VideoRef
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

type SavedVideo struct {
	ID       uuid.UUID `storm:"id"`
	VideoRef VideoRef  `storm:"unique"`
	Title    string
	FeedID   uuid.UUID
	Source   string
	SavedOn  time.Time
	Location string
	FileSize int64
}

type DownloadStatus int

const (
	StatusUnknown DownloadStatus = iota
	StatusNotDownloaded
	StatusDownloaded
	StatusMissing
)

func (ds DownloadStatus) String() string {
	switch ds {
	case StatusUnknown:
		return "Unknown"
	case StatusNotDownloaded:
		return "Not Downloaded"
	case StatusDownloaded:
		return "Downloaded"
	case StatusMissing:
		return "Missing"
	}
	return "Unknown"
}
