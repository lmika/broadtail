package jobs

import (
	"context"
	"github.com/google/uuid"
	"os"
	"sync"
)

type JobState int

const (
	StateQueued JobState = iota + 1
	StateRunning
	StateCancelling
	StateCompleted
	StateError
	StateCancelled
)

func (js JobState) String() string {
	switch js {
	case StateQueued:
		return "Queued"
	case StateRunning:
		return "Running"
	case StateCompleted:
		return "Done"
	case StateError:
		return "Error"
	case StateCancelling:
		return "Cancelling"
	case StateCancelled:
		return "Cancelled"
	}
	return "Unknown"
}

func (js JobState) Terminal() bool {
	switch js {
	case StateCompleted, StateCancelled, StateError:
		return true
	}
	return false
}

// Job is a handle to a running job
type Job struct {
	id   uuid.UUID
	task Task

	// state variables.  These are under the control of the mutex
	stateMutex sync.Mutex
	state      JobState
	err        error
	cancelFn   func()
	lastUpdate Update
	data       map[string]interface{}
}

func newJob(task Task) *Job {
	return &Job{
		id:         uuid.New(),
		task:       task,
		state:      StateQueued,
		stateMutex: sync.Mutex{},
		data:       make(map[string]interface{}),
	}
}

// exec executes the job.  This is assumed to be executed by a go routine
func (j *Job) exec(ctx context.Context, runContext RunContext) {
	execCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	j.setState(runContext, StateRunning, nil, cancelFn)

	err := j.task.Execute(execCtx, runContext)
	select {
	case <-execCtx.Done():
		j.setState(runContext, StateCancelled, execCtx.Err(), nil)
	default:
		if err != nil {
			j.setState(runContext, StateError, err, nil)
		} else {
			j.setState(runContext, StateCompleted, nil, nil)
		}
	}
}

func (j *Job) ID() uuid.UUID {
	return j.id
}

// Task returns the original task of the job
func (j *Job) Task() Task {
	return j.task
}

func (j *Job) State() JobState {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	return j.state
}

func (j *Job) Data() map[string]interface{} {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	dataCp := make(map[string]interface{})
	for k, v := range j.data {
		dataCp[k] = v
	}

	return dataCp
}

func (j *Job) GetData(key string) interface{} {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	return j.data[key]
}

func (j *Job) SetData(key string, value interface{}) {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	j.data[key] = value
}

// setState sets the job state in a safe manner
func (j *Job) setState(runContext RunContext, state JobState, err error, cancelFn func()) {
	var oldState JobState

	// Critical section
	j.stateMutex.Lock()

	oldState = j.state
	j.state = state
	j.err = err
	j.cancelFn = cancelFn

	j.stateMutex.Unlock()
	// End critical section

	if oldState != state {
		runContext.postStateChange(oldState, state)
	}
}

// setState sets the job state in a safe manner
func (j *Job) setLastUpdate(lastUpdate Update) {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	j.lastUpdate = lastUpdate
}

func (j *Job) LastUpdate() Update {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	return j.lastUpdate
}

// Cancel cancels a running job
func (j *Job) Cancel() {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	if j.state == StateCancelling {
		return
	}

	if j.cancelFn != nil {
		j.state = StateCancelling
		j.cancelFn()
	}
}

// Err returns the error from running the job
func (j *Job) Err() error {
	j.stateMutex.Lock()
	defer j.stateMutex.Unlock()

	return j.err
}

func (j *Job) Cleanup() {
	if tmpDir, ok := j.data["_tmp_dir"].(string); ok {
		os.RemoveAll(tmpDir)
	}
}
