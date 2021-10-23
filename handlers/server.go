package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/jobs"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/ujs"
	"github.com/lmika/broadtail/services/ytdownload"
	"github.com/pkg/errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Config struct {
	LibraryDir string

	TemplateFS fs.FS
	AssetFS    fs.FS
}

func Server(config Config) (http.Handler, error) {
	dispatcher := jobs.New()
	dispatcherJobRecorderAndCleanupHandler(dispatcher)

	tmpls, err := template.ParseFS(config.TemplateFS, "*.html")
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse templates")
	}

	ytdownloadService := ytdownload.New(ytdownload.Config{
		LibraryDir: config.LibraryDir,
	})

	ytdownloadHandlers := &youTubeDownloadHandlers{ytdownloadService: ytdownloadService}
	jobsHandlers := &jobsHandlers{}

	r := mux.NewRouter()
	r.Handle("/", indexHandler()).Methods("GET")
	r.Handle("/job/download/youtube", ytdownloadHandlers.CreateDownloadJob()).Methods("POST")
	r.Handle("/jobs/done", jobsHandlers.ClearDone()).Methods("DELETE")
	r.Handle("/jobs/{job_id}", jobsHandlers.Delete()).Methods("DELETE")
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(config.AssetFS))))

	handler := ujs.MethodRewriteHandler(r)
	handler = jobdispatcher.New(dispatcher).Use(handler)
	handler = render.New(tmpls).Use(handler)

	return handler, nil
}

func dispatcherJobRecorderAndCleanupHandler(dispatcher *jobs.Dispatcher) {
	sub := dispatcher.Subscribe()

	go func() {
		defer sub.Close()

		for event := range sub.Chan() {
			log.Printf("received event: %v", event)
			switch e := event.(type) {
			case jobs.StateTransitionSubscriptionEvent:
				if e.ToState.Terminal() {
					// TODO: save the done jobs
					doneJobs := dispatcher.ClearDone()
					log.Printf("Cleaned up %v jobs", len(doneJobs))
				}
			}
		}
	}()
}