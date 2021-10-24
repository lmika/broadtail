package stormstore

import (
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
)

type JobStore struct {
	db	*storm.DB
}

func NewJobStore(filename string) (*JobStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &JobStore{db: db}, nil
}

func (js *JobStore) Close() {
	js.db.Close()
}

func (js *JobStore) Job(id uuid.UUID) (job models.Job, err error) {
	err = js.db.One("ID", id, &job)
	return job, err
}

func (js *JobStore) List() (jobs []models.Job, err error) {
	err = js.db.Select().OrderBy("CreatedAt").Reverse().Limit(50).Find(&jobs)
	if err == storm.ErrNotFound {
		return []models.Job{}, nil
	}

	return jobs, err
}

func (js *JobStore) Save(job models.Job) error {
	return js.db.Save(&job)
}