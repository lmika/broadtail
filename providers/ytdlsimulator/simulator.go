package ytdlsimulator

import (
	"context"
	"fmt"
	"time"

	"github.com/lmika/broadtail/models"
)

type YoutubeDLSimulator struct{}

func New() YoutubeDLSimulator {
	return YoutubeDLSimulator{}
}

func (YoutubeDLSimulator) GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error) {
	return &models.Video{
		ExtID:        youtubeId,
		Title:        "Simulated video",
		Description:  "A simulated video with the external ID = " + youtubeId,
		ThumbnailURL: "https://www.example.com/",
		UploadedOn:   time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC),
	}, nil
}

func (YoutubeDLSimulator) DownloadVideo(ctx context.Context, options models.DownloadOptions, logline func(line string)) error {
	for i := 1; i <= 100; i++ {
		logline(fmt.Sprintf("[download] %d.0%% of 269.30MiB at 45.79KiB/s ETA 00:00", i))
		time.Sleep(1 * time.Second)
	}
	return nil
}
