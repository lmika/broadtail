package videomanager

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type VideoManager struct {
	dataDir    string
	videoStore VideoStore
}

func New(dataDir string, videoStore VideoStore) *VideoManager {
	return &VideoManager{
		dataDir:    dataDir,
		videoStore: videoStore,
	}
}

func (vm *VideoManager) List() ([]models.SavedVideo, error) {
	return vm.videoStore.ListRecent()
}

func (vm *VideoManager) Get(id uuid.UUID) (*models.SavedVideo, error) {
	return vm.videoStore.FindWithID(id)
}

//func (vm *VideoManager) DownloadStatus(extId string) (models.DownloadStatus, error) {
func (vm *VideoManager) DownloadStatus(videoRef models.VideoRef) (models.DownloadStatus, error) {
	video, err := vm.videoStore.FindWithExtID(videoRef.String())
	if err != nil {
		return models.StatusUnknown, err
	}

	if video == nil {
		return models.StatusNotDownloaded, nil
	}

	_, err = os.Stat(filepath.Join(vm.dataDir, video.Location))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return models.StatusMissing, nil
		} else {
			return models.StatusUnknown, err
		}
	}

	return models.StatusDownloaded, nil
}
