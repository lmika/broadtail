package jobs

import (
	"container/list"
	"context"
	"github.com/google/uuid"
	"sync"
)

type Dispatcher struct {
	//waitList	*list.List
	waitList chan *Job
	running  bool

	cancelFn  func()
	closeChan chan struct{}

	jobStore *SliceJobStore

	subManagementChan chan subMgmtEvent
}

func New(options ...Option) *Dispatcher {
	dispatcher := &Dispatcher{
		//waitList: list.New(),
		waitList:          make(chan *Job, 50),
		running:           true,
		subManagementChan: make(chan subMgmtEvent, 50),
		jobStore:          NewSliceJobStore(),
	}

	for _, option := range options {
		option(dispatcher)
	}

	if dispatcher.running {
		dispatcher.start()
	}

	return dispatcher
}

// List returns the list of active jobs
func (d *Dispatcher) List() []*Job {
	return d.jobStore.List()
}

// Job returns the job with the given ID.
func (d *Dispatcher) Job(id uuid.UUID) *Job {
	return d.jobStore.Find(id)
}

// ClearDone removes all the jobs that are not running
func (d *Dispatcher) ClearDone() {
	for _, clearedJob := range d.jobStore.ClearDone() {
		clearedJob.Cleanup()
	}
}

// Enqueue adds a new job to the wait queue.  If there is something that can
// execute the job immediately, it will be started.
func (d *Dispatcher) Enqueue(task Task) *Job {
	job := newJob(task)
	d.jobStore.Add(job)

	//d.waitList.PushBack(job)
	d.waitList <- job
	return job
}

// Subscribe starts a new subscription for job updates.
func (d *Dispatcher) Subscribe() *Subscription {
	respChan := make(chan *Subscription)
	d.subManagementChan <- subMgmtNewSubscription{respChan}
	return <- respChan
}

// Close cancels the running tasks and waits for the dispatcher to stop
func (d *Dispatcher) Close() {
	if d.cancelFn == nil {
		return
	}

	d.cancelFn()
	<-d.closeChan
}

// drain will execute all the jobs in the foreground.  This is used for testing
// purposes.
//func (d *Dispatcher) drain() {
//	for e := d.waitList.Front(); e != nil; e = e.Next() {
//		j := d.waitList.Remove(e).(*Job)
//		d.startJob(context.Background(), j)
//	}
//}

func (d *Dispatcher) start() {
	runCtx, cancelFn := context.WithCancel(context.Background())
	d.cancelFn = cancelFn
	d.closeChan = make(chan struct{})

	go func() {
		d.subscriptionManager(runCtx, d.subManagementChan)
		close(d.subManagementChan)
	}()

	go func() {
		d.loop(runCtx)
		close(d.closeChan)
	}()
}

func (d *Dispatcher) subscriptionManager(ctx context.Context, eventChan chan subMgmtEvent) {
	subList := list.New()

	for {
		select {
		case <- ctx.Done():
			return
		case event := <- eventChan:
			switch e := event.(type) {
			case subMgmtNewSubscription:
				newSub := &Subscription{c: make(chan Update, 5)}
				newSub.elem = subList.PushBack(newSub)
				newSub.closeFn = func() {
					eventChan <- subMgmtUnsubscribe{newSub}
				}
				e.retChan <- newSub
			case subMgmtUnsubscribe:
				subList.Remove(e.sub.elem)
			case subMgmtPublish:
				for se := subList.Front(); se != nil; se = se.Next() {
					se.Value.(*Subscription).c <- e.update
				}
			}
		}
	}
}

//func (d *Dispatcher) publish(update Update) {
//	d.subManagementChan <- subMgmtPublish{update}
//}

func (d *Dispatcher) loop(ctx context.Context) {
	const workers int = 4

	dispatchChan := make(chan *Job)
	var waitGroup sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()

			for job := range dispatchChan {
				d.startJob(ctx, job)
			}
		}()
	}

mainFor:
	for {
		select {
		case job := <-d.waitList:
			dispatchChan <- job
		case <-ctx.Done():
			break mainFor
		}
	}

	// Close everything
	close(dispatchChan)
	waitGroup.Wait()
}

func (d *Dispatcher) startJob(ctx context.Context, job *Job) {
	job.exec(ctx, &jobRunContext{
		job: job,
		subManagementChan: d.subManagementChan,
	})
}

type Option func(disp *Dispatcher)

func StartPaused() func(disp *Dispatcher) {
	return func(disp *Dispatcher) { disp.running = false }
}

type subMgmtEvent interface{}

type subMgmtNewSubscription struct {
	retChan chan *Subscription
}

type subMgmtPublish struct {
	update Update
}

type subMgmtUnsubscribe struct {
	sub *Subscription
}
