package handlers

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/rssfetcher"
	"github.com/lmika/broadtail/middleware/sessions"
	"github.com/lmika/broadtail/middleware/ujs"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/lmika/broadtail/providers/stormstore"
	"github.com/lmika/broadtail/providers/youtubedl"
	"github.com/lmika/broadtail/providers/ytdlsimulator"
	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/broadtail/services/ytdownload"
	"github.com/pkg/errors"
)

type Config struct {
	LibraryDir          string
	JobDataFile         string
	FeedsDataFile       string
	YTDownloadSimulator bool

	TemplateFS fs.FS
	AssetFS    fs.FS
}

func Server(config Config) (handler http.Handler, closeFn func(), err error) {
	dispatcher := jobs.New()

	jobStore, err := stormstore.NewJobStore(config.JobDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open job store")
	}
	feedsStore, err := stormstore.NewFeedStore(config.FeedsDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open feeds store")
	}
	rssFetcher := rssfetcher.New()

	var youtubedlProvider ytdownload.Provider
	if config.YTDownloadSimulator {
		log.Println("Using youtuble-dl simulator")
		youtubedlProvider = ytdlsimulator.New()
	} else {
		youtubedlProvider, err = youtubedl.New()
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot instantiate youtube-dl provider")
		}
	}

	ytdownloadService := ytdownload.New(ytdownload.Config{LibraryDir: config.LibraryDir}, youtubedlProvider, feedsStore)
	feedsManager := feedsmanager.New(feedsStore, rssFetcher)
	jobsManager := jobsmanager.New(dispatcher, jobStore)
	go jobsManager.Start()

	indexHandlers := &indexHandlers{jobsManager: jobsManager}
	ytdownloadHandlers := &youTubeDownloadHandlers{ytdownloadService: ytdownloadService, jobsManager: jobsManager}
	jobsHandlers := &jobsHandlers{jobsManager: jobsManager}
	feedsHandlers := &feedsHandler{feedsManager: feedsManager}
	configHandlers := &configHandler{}

	r := mux.NewRouter()
	r.Handle("/", indexHandlers.Index()).Methods("GET")
	r.Handle("/video/details", ytdownloadHandlers.ShowDetails()).Methods("GET")
	r.Handle("/job/download/youtube", ytdownloadHandlers.CreateDownloadJob()).Methods("POST")

	r.Handle("/jobs", jobsHandlers.List()).Methods("GET")
	r.Handle("/jobs/done", jobsHandlers.ClearDone()).Methods("DELETE")
	r.Handle("/jobs/{job_id}", jobsHandlers.Show()).Methods("GET")
	r.Handle("/jobs/{job_id}", jobsHandlers.Delete()).Methods("DELETE")

	r.Handle("/config", configHandlers.Show()).Methods("GET")

	feedsHandlers.Routes(r)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(config.AssetFS))))

	handler = ujs.MethodRewriteHandler(r)
	handler = jobdispatcher.New(dispatcher).Use(handler)
	handler = sessions.Use(handler)
	handler = render.New(config.TemplateFS).Use(handler)

	closeFn = func() {
		jobStore.Close()
		feedsStore.Close()
	}

	return handler, closeFn, nil
}
