package ytdownload

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
)

type Config struct {
	LibraryDir string
}

type Service struct {
	config     Config
	provider   DownloadProvider
	feedStore  FeedStore
	videoStore VideoStore
}

func New(config Config, provider DownloadProvider, feedStore FeedStore, videoStore VideoStore) *Service {
	return &Service{
		config:     config,
		provider:   provider,
		feedStore:  feedStore,
		videoStore: videoStore,
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error) {
	return s.provider.GetVideoMetadata(ctx, youtubeId)
}

func (s *Service) NewYoutubeDownloadTask(ctx context.Context, youtubeId string, feedIDConfig uuid.UUID) (jobs.Task, error) {
	targetDir := s.config.LibraryDir

	var sourceFeed *models.Feed
	if feedIDConfig != uuid.Nil {
		feed, err := s.feedStore.Get(ctx, feedIDConfig)
		if err != nil {
			return nil, err
		}

		sourceFeed = &feed
		targetDir = filepath.Join(s.config.LibraryDir, feed.TargetDir)
	}

	task := &YoutubeDownloadTask{
		DownloadProvider: s.provider,
		Feed:             sourceFeed,
		YoutubeId:        youtubeId,
		TargetDir:        targetDir,
		VideoStore:       s.videoStore,
	}
	task.Init()
	return task, nil
}
