package jobs

import (
	"github.com/google/uuid"
	"sync"
)

type SliceJobStore struct {
	mutex	sync.Mutex
	jobs		[]*Job
}

func NewSliceJobStore() *SliceJobStore {
	return &SliceJobStore{
		mutex: sync.Mutex{},
		jobs: make([]*Job, 0),
	}
}

func (js *SliceJobStore) Add(job *Job) {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	js.jobs = append(js.jobs, job)
}

func (js *SliceJobStore) List() []*Job {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	return js.jobs
}

func (js *SliceJobStore) ClearDone() []*Job {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	newJobs := make([]*Job, 0)
	deletedJobs := make([]*Job, 0)
	for _, j := range js.jobs {
		if !j.state.Terminal() {
			newJobs = append(newJobs, j)
		} else {
			deletedJobs = append(deletedJobs, j)
		}
	}

	js.jobs = newJobs
	return deletedJobs
}

func (js *SliceJobStore) Find(id uuid.UUID) *Job {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	for _, j := range js.jobs {
		if j.id == id {
			return j
		}
	}
	return nil
}
