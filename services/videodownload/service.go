package videodownload

import (
	"context"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/broadtail/services/ytdownload"
	"github.com/pkg/errors"
)

type Service struct {
	ytdl        *ytdownload.Service
	jobsManager *jobsmanager.JobsManager
}

func NewService(ytdl *ytdownload.Service, jobsManager *jobsmanager.JobsManager) *Service {
	return &Service{
		ytdl:        ytdl,
		jobsManager: jobsManager,
	}
}

// TODO: remove use of feedID
func (s *Service) QueueForDownload(ctx context.Context, videoRef models.VideoRef, feedID uuid.UUID) error {
	if videoRef.Source != models.YoutubeVideoRefSource {
		return errors.Errorf("unrecognised video ref source: %v", videoRef.Source)
	}

	// TODO: handle support for different video type
	task, err := s.ytdl.NewYoutubeDownloadTask(ctx, videoRef.ID, feedID)
	if err != nil {
		return err
	}

	s.jobsManager.Dispatcher().Enqueue(task)
	return nil
}
