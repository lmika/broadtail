package models

import (
	"github.com/google/uuid"
	"github.com/lmika/broadtail/providers/jobs"
	"time"
)

type JobUpdate struct {
	Message string
}

type Job struct {
	ID          uuid.UUID `storm:"unique"`
	CreatedAt   time.Time `storm:"index"`
	CompletedAt time.Time
	Name        string
	State       jobs.JobState
	Updates     []JobUpdate
	Error       string
}

func (j *Job) SetLastUpdate(update JobUpdate) {
	j.Updates = []JobUpdate{update}
}

func (j Job) LastMessage() string {
	if j.Error != "" {
		return j.Error
	} else if len(j.Updates) == 0 {
		return ""
	}

	return j.Updates[len(j.Updates)-1].Message
}
