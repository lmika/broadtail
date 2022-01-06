package models

import (
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lmika/broadtail/providers/jobs"
)

type JobUpdate struct {
	Message string
}

type Job struct {
	ID          uuid.UUID `storm:"unique"`
	CreatedAt   time.Time `storm:"index"`
	CompletedAt time.Time
	Name        string
	VideoExtID  string
	VideoTitle  string
	State       jobs.JobState
	Progress    Progress
	Updates     []JobUpdate
	Error       string
}

func (j *Job) SetLastUpdate(update JobUpdate) {
	j.Updates = []JobUpdate{update}
	if progress, ok := ParseProgress(update.Message); ok {
		j.Progress = progress
	}
}

func (j Job) LastMessage() string {
	if j.Error != "" {
		return j.Error
	} else if len(j.Updates) == 0 {
		return ""
	}

	return j.Updates[len(j.Updates)-1].Message
}

type Progress struct {
	Percent float64
	Size    string
	Rate    string
	ETA     string
}

func ParseProgress(message string) (Progress, bool) {
	groups := progressRegexp.FindStringSubmatch(message)
	if len(groups) != 5 {
		return Progress{}, false
	}

	percentFloat, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return Progress{}, false
	}

	return Progress{
		Percent: percentFloat,
		Size:    groups[2],
		Rate:    groups[3],
		ETA:     groups[4],
	}, false
}

var progressRegexp = regexp.MustCompile(`\[download\] ([0-9.]+)% of ([0-9A-Za-z.]+) at ([0-9A-Za-z.]+)/s ETA ([0-9:.]+)`)
