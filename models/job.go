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
	//Updates     []JobUpdate
	LastUpdate JobUpdate
	Messages   []string
	Error      string
}

//func (j *Job) SetLastMessage(update string) {
//	j.Messages = []string{update}
//}

//func (j Job) LastMessage() string {
//}

func (j Job) LastMessage() string {
	if j.Error != "" {
		return j.Error
	} else if len(j.Messages) == 0 {
		return ""
	}

	return j.Messages[len(j.Messages)-1]
}

//type progress struct {
//	Percent float64
//	Size    string
//	Rate    string
//	ETA     string
//}

type JobUpdate struct {
	Summary string
	Percent float64
}

/*
func ParseJobUpdate(message string) JobUpdate {
	p, _ := parseProgress(message)
	return JobUpdate{
		Message: message,
		Percent: p.Percent,
	}
}

func parseProgress(message string) (progress, bool) {
	log.Printf("Parsing progress: '%v'", message)
	groups := progressRegexp.FindStringSubmatch(message)
	if len(groups) != 5 {
		log.Println("bad progress: ", len(groups))
		return progress{}, false
	}

	percentFloat, err := strconv.ParseFloat(groups[1], 64)
	if err != nil {
		return progress{}, false
	}

	return progress{
		Percent: percentFloat,
		Size:    groups[2],
		Rate:    groups[3],
		ETA:     groups[4],
	}, false
}

// [download] 2.1% of 86.31MiB at 84.91KiB/s ETA 16:59
var progressRegexp = regexp.MustCompile(`\[download\]\s+([0-9.]+)% of ([0-9A-Za-z.]+) at ([0-9A-Za-z.]+)/s ETA ([0-9:.]+)`)
*/
