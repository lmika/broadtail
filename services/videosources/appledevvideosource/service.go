package appledevvideosource

import (
	"context"
	"fmt"
	"github.com/lmika/broadtail/models"
	"net/url"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	return &models.Video{
		VideoRef:    videoRef,
		Title:       "Some title from apple",
		Description: "Some description from apple",
		Duration:    0,
	}, nil
}

func (s *Service) DownloadVideo(ctx context.Context, videoRef models.VideoRef, options models.DownloadOptions, logline func(line string)) (outputFilename string, err error) {
	panic("implement me")
}

func (s *Service) GetVideoURL(videoRef models.VideoRef) string {
	videoSet, video, _ := strings.Cut(videoRef.ID, ".")
	return fmt.Sprintf("https://developer.apple.com/videos/play/%s/%s/", url.PathEscape(videoSet), url.PathEscape(video))
}
