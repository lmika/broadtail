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
	config    Config
	provider  Provider
	feedStore FeedStore
}

func New(config Config, provider Provider, feedStore FeedStore) *Service {
	return &Service{
		config:    config,
		provider:  provider,
		feedStore: feedStore,
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, youtubeId string) (*models.Video, error) {
	return s.provider.GetVideoMetadata(ctx, youtubeId)
}

func (s *Service) NewYoutubeDownloadTask(ctx context.Context, youtubeId string, feedIDConfig uuid.UUID) (jobs.Task, error) {
	targetDir := s.config.LibraryDir
	if feedIDConfig != uuid.Nil {
		feed, err := s.feedStore.Get(ctx, feedIDConfig)
		if err != nil {
			return nil, err
		}

		targetDir = filepath.Join(s.config.LibraryDir, feed.TargetDir)
	}

	task := &YoutubeDownloadTask{
		Provider:  s.provider,
		YoutubeId: youtubeId,
		TargetDir: targetDir,
	}
	task.Init()
	return task, nil
}
