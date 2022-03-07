package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/providers/jobs"
)

type Job struct {
	ID          uuid.UUID `storm:"unique"`
	CreatedAt   time.Time `storm:"index"`
	CompletedAt time.Time
	Name        string
	VideoExtID  string
	VideoTitle  string
	State       jobs.JobState
	LastUpdate  JobUpdate
	Messages    []string
	Error       string
}

func (j Job) LastMessage() string {
	if j.Error != "" {
		return j.Error
	} else if len(j.Messages) == 0 {
		return ""
	}

	return j.Messages[len(j.Messages)-1]
}

type JobUpdate struct {
	Summary string
	Percent float64
}
