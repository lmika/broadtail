package jobsmanager

import (
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/pkg/errors"
)

type JobStore interface {
	Job(id uuid.UUID) (models.Job, error)
	List() ([]models.Job, error)
	Save(job models.Job) error
}

type JobsManager struct {
	dispatcher *jobs.Dispatcher
	jobStore   JobStore
}

func New(dispatcher *jobs.Dispatcher, jobStore JobStore) *JobsManager {
	return &JobsManager{dispatcher: dispatcher, jobStore: jobStore}
}

func (jm *JobsManager) Dispatcher() *jobs.Dispatcher {
	return jm.dispatcher
}

func (jm *JobsManager) Start() {
	sub := jm.dispatcher.Subscribe()

	go func() {
		defer sub.Close()

		for event := range sub.Chan() {
			switch e := event.(type) {
			case jobs.StateTransitionSubscriptionEvent:
				if e.ToState.Terminal() {
					// Save this job
					jobToSave := jm.toJob(e.Job, true)
					jobToSave.CompletedAt = time.Now()
					if err := jm.jobStore.Save(jobToSave); err != nil {
						log.Printf("warn: unable to save job: %v", err)
					}

					jm.dispatcher.ClearDone()
				}
			}
		}
	}()
}

func (jm *JobsManager) RecentJobs() []models.Job {
	dispatcherJobs := jm.dispatcher.List()
	js := make([]models.Job, len(dispatcherJobs))
	for i, j := range dispatcherJobs {
		js[i] = jm.toJob(j, false)
	}

	sort.Slice(js, func(i, j int) bool {
		return js[i].CreatedAt.Before(js[j].CreatedAt)
	})

	return js
}

func (jm *JobsManager) HistoricalJobs() ([]models.Job, error) {
	historical, err := jm.jobStore.List()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get historical jobs")
	}
	return historical, nil
}

func (jm *JobsManager) Job(id uuid.UUID) (models.Job, error) {
	runningJob := jm.dispatcher.Job(id)
	if runningJob != nil {
		return jm.toJob(runningJob, true), nil
	}

	return jm.jobStore.Job(id)
}

func (jm *JobsManager) toJob(job *jobs.Job, updateHistory bool) models.Job {
	modelJob := models.Job{
		ID:        job.ID(),
		Name:      job.Task().String(),
		CreatedAt: job.CreatedAt(),
		State:     job.State(),
	}
	if videoDownloadTask, isVideoDownloadTask := job.Task().(VideoDownloadTask); isVideoDownloadTask {
		modelJob.VideoExtID = videoDownloadTask.VideoExtID()
		modelJob.VideoTitle = videoDownloadTask.VideoTitle()
	}

	if err := job.Err(); err != nil {
		modelJob.Error = err.Error()
	}

	if updateHistory {
		history := job.UpdateHistory()
		modelJob.Updates = make([]models.JobUpdate, len(history))
		for i, h := range history {
			modelJob.Updates[i] = models.JobUpdate{Message: h.Status}
		}
	} else {
		modelJob.Updates = []models.JobUpdate{models.ParseJobUpdate(job.LastUpdate().Status)}
	}

	return modelJob
}
