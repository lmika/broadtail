package handlers

import (
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/ujs"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/lmika/broadtail/providers/stormstore"
	"github.com/lmika/broadtail/providers/youtubedl"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/broadtail/services/ytdownload"
	"github.com/pkg/errors"
	"io/fs"
	"net/http"
)

type Config struct {
	LibraryDir     string
	CacheTemplates bool
	JobDataFile        string

	TemplateFS fs.FS
	AssetFS    fs.FS
}

func Server(config Config) (handler http.Handler, closeFn func(), err error) {
	dispatcher := jobs.New()
	jobStore, err := stormstore.NewJobStore(config.JobDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open job store")
	}

	youtubedlProvider := youtubedl.New()

	ytdownloadService := ytdownload.New(ytdownload.Config{LibraryDir: config.LibraryDir}, youtubedlProvider)
	jobsManager := jobsmanager.New(dispatcher, jobStore)
	go jobsManager.Start()

	indexHandlers := &indexHandlers{jobsManager: jobsManager}
	ytdownloadHandlers := &youTubeDownloadHandlers{ytdownloadService: ytdownloadService, jobsManager: jobsManager}
	jobsHandlers := &jobsHandlers{jobsManager: jobsManager}

	r := mux.NewRouter()
	r.Handle("/", indexHandlers.Index()).Methods("GET")
	r.Handle("/video/details", ytdownloadHandlers.ShowDetails()).Methods("GET")
	r.Handle("/job/download/youtube", ytdownloadHandlers.CreateDownloadJob()).Methods("POST")

	r.Handle("/jobs", jobsHandlers.List()).Methods("GET")
	r.Handle("/jobs/done", jobsHandlers.ClearDone()).Methods("DELETE")
	r.Handle("/jobs/{job_id}", jobsHandlers.Show()).Methods("GET")
	r.Handle("/jobs/{job_id}", jobsHandlers.Delete()).Methods("DELETE")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(config.AssetFS))))

	handler = ujs.MethodRewriteHandler(r)
	handler = jobdispatcher.New(dispatcher).Use(handler)
	handler = render.New(config.TemplateFS, config.CacheTemplates).Use(handler)

	closeFn = func() {
		jobStore.Close()
	}

	return handler, closeFn, nil
}
