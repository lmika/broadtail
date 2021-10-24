package jobs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
)

// Task is something that a job is to perform
type Task interface {
	String() string

	// Execute executes the task with the given context.
	Execute(ctx context.Context, runContext RunContext) error
}

type DisplayableTask interface {
	PreferredTemplate() string
}

type RunContext interface {
	PostUpdate(update Update)
	TempFile(pattern string) (*os.File, error)
	Set(key string, value interface{})

	postStateChange(fromState, toState JobState)
}

func PostUpdatef(runContext RunContext, msg string, args ...interface{}) {
	runContext.PostUpdate(Update{fmt.Sprintf(msg, args...)})
}

type jobRunContext struct {
	job *Job
	subManagementChan chan subMgmtEvent
}

func (rc *jobRunContext) PostUpdate(update Update) {
	rc.job.postUpdate(update)
	rc.subManagementChan <- subMgmtPublish{UpdateSubscriptionEvent{rc.job, update}}
}

func (rc *jobRunContext) postStateChange(fromState, toState JobState) {
	rc.subManagementChan <- subMgmtPublish{StateTransitionSubscriptionEvent{rc.job, fromState, toState}}
}

func (rc *jobRunContext) TempFile(pattern string) (*os.File, error) {
	var tmpDir string
	var ok bool

	tmpDir, ok = rc.job.GetData("_tmp_dir").(string)
	if !ok {
		newTmpDir, err := ioutil.TempDir("", "job-workspace")
		if err != nil {
			return nil, err
		}

		rc.job.SetData("_tmp_dir", newTmpDir)
	}

	return ioutil.TempFile(tmpDir, pattern)
}

func (rc *jobRunContext) Set(key string, value interface{}) {
	rc.job.SetData(key, value)
}