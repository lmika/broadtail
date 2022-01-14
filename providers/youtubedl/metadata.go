package youtubedl

import (
	"time"
)

type metadataJson struct {
	UploadDateStr string `json:"upload_date"`
	Title         string `json:"title"`
	ChannelID     string `json:"channel_id"`
	Channel       string `json:"channel"`
	Description   string `json:"description"`
	ThumbnailURL  string `json:"thumbnail"`
	Duration      int    `json:"duration"`
}

func (r metadataJson) UploadDate() (time.Time, error) {
	return time.Parse("20060102", r.UploadDateStr)
}
