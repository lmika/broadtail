package videodownload

import (
	"context"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/broadtail/services/videosources"
)

type Config struct {
	LibraryDir          string
	LibraryOwner        string
	VideoSourcesService *videosources.Service
	VideoStore          VideoStore
	VideoDownloadHooks  VideoDownloadHooks
	FeedStore           FeedStore
	FeedItemStore       FeedItemStore
	JobsManager         *jobsmanager.JobsManager
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	return &Service{config: config}
}

func (s *Service) GetVideoMetadata(ctx context.Context, videoRef models.VideoRef) (*models.Video, error) {
	videoSource, err := s.config.VideoSourcesService.SourceProvider(videoRef)
	if err != nil {
		return nil, err
	}

	return videoSource.GetVideoMetadata(ctx, videoRef)
}

// TODO: remove use of feedID
func (s *Service) QueueForDownload(ctx context.Context, videoRef models.VideoRef, feed *models.Feed) error {
	videoSource, err := s.config.VideoSourcesService.SourceProvider(videoRef)
	if err != nil {
		return err
	}

	var targetDir = ""

	var videoFeed *models.Feed
	if feed != nil {
		videoFeed = feed
	} else {
		if feedItem, err := s.config.FeedItemStore.GetByVideoRef(ctx, videoRef); err == nil {
			if feed, err := s.config.FeedStore.Get(ctx, feedItem.FeedID); err == nil {
				videoFeed = &feed
				targetDir = feed.TargetDir
			}
		}
	}

	task := &YoutubeDownloadTask{
		VideoSource:       videoSource,
		Feed:              videoFeed,
		VideoRef:          videoRef,
		LibraryDir:        s.config.LibraryDir,
		TargetDir:         targetDir,
		TargetOwner:       s.config.LibraryOwner,
		VideoStore:        s.config.VideoStore,
		VideoDownloadHook: s.config.VideoDownloadHooks,
	}
	task.Init()

	s.config.JobsManager.Dispatcher().Enqueue(task)
	return nil
}
