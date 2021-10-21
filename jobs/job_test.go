package jobs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJob_Normal(t *testing.T) {
	task := &waitTask{ waitTime: 1 * time.Second }

	disp := New()
	job := disp.Enqueue(task)

	// FIXME: wait for subscriptions
	time.Sleep(2 * time.Second)

	assert.Equal(t, true, task.finished)
	assert.Equal(t, StateCompleted, job.State())
	assert.Nil(t, job.Err())
}

func TestJob_Cancel(t *testing.T) {
	task := &waitTask{ waitTime: 5 * time.Second }

	disp := New()
	job := disp.Enqueue(task)

	// FIXBE
	time.Sleep(2 * time.Second)
	job.Cancel()

	time.Sleep(2 * time.Second)

	assert.Equal(t, true, task.cancelled)
	assert.Equal(t, StateCancelled, job.State())
	assert.NotNil(t, job.Err())
}


type waitTask struct {
	waitTime 	time.Duration
	finished bool
	cancelled bool
}

func (s *waitTask) Execute(ctx context.Context) error {
	timer := time.NewTimer(s.waitTime)

	select {
	case <- timer.C:
		s.finished = true
	case <- ctx.Done():
		s.cancelled = true
	}
	return nil
}