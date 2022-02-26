package stormstore

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type VideoStore struct {
	db *storm.DB
}

func (vs *VideoStore) DeleteWithExtID(extId string) error {
	if err := vs.db.Select(q.Eq("ExtID", extId)).Delete(&models.SavedVideo{}); err != nil {
		if err == storm.ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (vs *VideoStore) FindWithID(id uuid.UUID) (*models.SavedVideo, error) {
	var savedVideo models.SavedVideo
	if err := vs.db.One("ID", id, &savedVideo); err != nil {
		//if err == storm.ErrNotFound {
		//	return nil, nil
		//}
		return nil, err
	}

	return &savedVideo, nil
}

func (vs *VideoStore) FindWithExtID(extId string) (*models.SavedVideo, error) {
	var savedVideo models.SavedVideo
	if err := vs.db.Select(q.Eq("ExtID", extId)).First(&savedVideo); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &savedVideo, nil
}

func NewVideoStore(filename string) (*VideoStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &VideoStore{db: db}, nil
}

func (vs *VideoStore) Close() {
	vs.db.Close()
}

func (vs *VideoStore) Save(video models.SavedVideo) error {
	return vs.db.Save(&video)
}

func (vs *VideoStore) ListRecent() (videos []models.SavedVideo, err error) {
	err = vs.db.Select().OrderBy("SavedOn").Reverse().Limit(50).Find(&videos)
	if err == storm.ErrNotFound {
		return []models.SavedVideo{}, nil
	}

	return videos, err
}
