package ytdownload

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
)

type Config struct {
	LibraryDir   string
	LibraryOwner string
}

type Service struct {
	config            Config
	provider          DownloadProvider
	feedStore         FeedStore
	videoStore        VideoStore
	videoDownloadHook VideoDownloadHooks
}

func New(config Config, provider DownloadProvider, feedStore FeedStore, videoStore VideoStore, videoDownloadHook VideoDownloadHooks) *Service {
	return &Service{
		config:            config,
		provider:          provider,
		feedStore:         feedStore,
		videoStore:        videoStore,
		videoDownloadHook: videoDownloadHook,
	}
}

func (s *Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	return s.provider.GetVideoMetadata(ctx, videoRef)
}

func (s *Service) NewYoutubeDownloadTask(ctx context.Context, videoRef models.VideoRef, feedIDConfig uuid.UUID) (jobs.Task, error) {
	panic("DELETE ME")
	//targetDir := s.config.LibraryDir
	//
	//var sourceFeed *models.Feed
	//if feedIDConfig != uuid.Nil {
	//	feed, err := s.feedStore.Get(ctx, feedIDConfig)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	sourceFeed = &feed
	//	targetDir = filepath.Join(s.config.LibraryDir, feed.TargetDir)
	//}

	//task := &YoutubeDownloadTask{
	//	DownloadProvider:  s.provider,
	//	Feed:              sourceFeed,
	//	VideoRef:          videoRef,
	//	TargetDir:         targetDir,
	//	TargetOwner:       s.config.LibraryOwner,
	//	VideoStore:        s.videoStore,
	//	VideoDownloadHook: s.videoDownloadHook,
	//}
	//task.Init()
	//return task, nil
}
