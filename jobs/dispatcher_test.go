package jobs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDispatcher_Enqueue(t *testing.T) {
	//disp := New(StartPaused())
	//task := &simpleTask{}
	//
	//job := disp.Enqueue(task)
	//assert.NotEqual(t, uuid.Nil, job.ID())
	//assert.Equal(t, StateQueued, job.State())
	//assert.Equal(t, task, job.Task())
	//
	//// Launch job
	//disp.drain()
	//
	//// Check job finished
	//assert.Equal(t, true, task.wasRun)
	//assert.Equal(t, StateCompleted, job.State())
}

func TestDispatcher_Close(t *testing.T) {
	task := &waitTask{ waitTime: 5 * time.Second }

	disp := New()
	job := disp.Enqueue(task)

	// FIXBE
	time.Sleep(2 * time.Second)
	disp.Close()

	time.Sleep(2 * time.Second)

	assert.Equal(t, true, task.cancelled)
	assert.Equal(t, StateCancelled, job.State())
	assert.NotNil(t, job.Err())
}


type simpleTask struct {
	wasRun	bool
}

func (s *simpleTask) Execute(ctx context.Context) error {
	s.wasRun = true
	return nil
}



