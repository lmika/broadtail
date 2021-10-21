package jobdispatcher

import (
	"context"
	"github.com/lmika/broadtail/jobs"
	"net/http"
)

type JobDispatcher struct {
	Dispatcher *jobs.Dispatcher
}

func New(dispatcher *jobs.Dispatcher) *JobDispatcher {
	return &JobDispatcher{dispatcher}
}

func (jd *JobDispatcher) Use(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), jobDispatcherContextKey, jd)))
	})
}

func FromContext(ctx context.Context) *JobDispatcher {
	rc, ok := ctx.Value(jobDispatcherContextKey).(*JobDispatcher)
	if !ok {
		return nil
	}
	return rc
}

type jobDispatcherContext struct {
	jobDispatcher *JobDispatcher
}

type jobDispatcherContextKeyType struct{}

var jobDispatcherContextKey = jobDispatcherContextKeyType{}
