package simulatorvideosource

import (
	"context"
	"fmt"
	"time"

	"github.com/lmika/broadtail/models"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

// GetVideoMetadata returns the metadata details for a sourceID.
func (Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	return &models.Video{
		ExtID:        videoRef.String(),
		Title:        "Simulated video",
		ChannelID:    "chan123",
		ChannelName:  "Simulated channel",
		Description:  "A simulated video with the external ID = " + videoRef.String(),
		ThumbnailURL: "https://www.example.com/",
		UploadedOn:   time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
	}, nil
}

// BuildDownloadTask returns a new task that will download the video.
func (Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line string)) (outputFilename string, err error) {
	for i := 1; i <= 100; i++ {
		logline(fmt.Sprintf("[download] %d.0%% of 269.30MiB at 45.79KiB/s ETA 00:00", i))
		time.Sleep(1 * time.Second)
	}
	return "", nil
}
