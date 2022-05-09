package videomanager

import (
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type VideoStore interface {
	ListRecent() ([]models.SavedVideo, error)
	FindWithID(id uuid.UUID) (*models.SavedVideo, error)
	FindWithExtID(videoRef models.VideoRef) (*models.SavedVideo, error)
}
