package handlers

import (
	"context"
	"fmt"
	"github.com/lmika/broadtail/services/videomanager"
	"github.com/mergestat/timediff"
	"github.com/robfig/cron"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/rssfetcher"
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
	VideoDataFile       string
	FeedsDataFile       string
	FeedItemsDataFile   string
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
	videoStore, err := stormstore.NewVideoStore(config.VideoDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open video store")
	}
	feedsStore, err := stormstore.NewFeedStore(config.FeedsDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open feeds store")
	}
	feedItemStore, err := stormstore.NewFeedItemStore(config.FeedItemsDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open feeds store")
	}
	rssFetcher := rssfetcher.New()

	var youtubedlProvider ytdownload.DownloadProvider
	if config.YTDownloadSimulator {
		log.Println("Using youtuble-dl simulator")
		youtubedlProvider = ytdlsimulator.New()
	} else {
		youtubedlProvider, err = youtubedl.New()
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot instantiate youtube-dl provider")
		}
	}

	ytdownloadService := ytdownload.New(ytdownload.Config{LibraryDir: config.LibraryDir}, youtubedlProvider, feedsStore, videoStore)
	feedsManager := feedsmanager.New(feedsStore, feedItemStore, rssFetcher)
	jobsManager := jobsmanager.New(dispatcher, jobStore)
	videoManager := videomanager.New(config.LibraryDir, videoStore)
	go jobsManager.Start()

	// Schedule updates every 15 minutes
	c := cron.New()
	if err := c.AddFunc("@every 15m", func() {
		if err := feedsManager.UpdateAllFeeds(context.Background()); err != nil {
			log.Printf("error updating all feeds: %v", err)
		}
	}); err != nil {
		return nil, nil, errors.Wrap(err, "invalid feed update schedule")
	}
	c.Start()

	indexHandlers := &indexHandlers{jobsManager: jobsManager}
	ytdownloadHandlers := &youTubeDownloadHandlers{ytdownloadService: ytdownloadService, jobsManager: jobsManager}
	detailsHandler := &detailsHandler{ytdownloadService: ytdownloadService, videoManager: videoManager}
	videoHandler := &videoHandlers{videoManager: videoManager}
	jobsHandlers := &jobsHandlers{jobsManager: jobsManager}
	feedsHandlers := &feedsHandler{feedsManager: feedsManager}

	r := mux.NewRouter()
	r.Handle("/", indexHandlers.Index()).Methods("GET")
	r.Handle("/job/download/youtube", ytdownloadHandlers.CreateDownloadJob()).Methods("POST")

	r.Handle("/quicklook", detailsHandler.QuickLook()).Methods("GET")
	r.Handle("/details/video/{video_id}", detailsHandler.VideoDetails()).Methods("GET")

	r.Handle("/videos", videoHandler.List()).Methods("GET")

	r.Handle("/jobs", jobsHandlers.List()).Methods("GET")
	r.Handle("/jobs/done", jobsHandlers.ClearDone()).Methods("DELETE")
	r.Handle("/jobs/{job_id}", jobsHandlers.Show()).Methods("GET")
	r.Handle("/jobs/{job_id}", jobsHandlers.Delete()).Methods("DELETE")

	feedsHandlers.Routes(r)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(config.AssetFS))))

	handler = ujs.MethodRewriteHandler(r)
	handler = jobdispatcher.New(dispatcher).Use(handler)
	handler = render.New(
		config.TemplateFS,
		render.WithFuncs(template.FuncMap{
			"formatTime": func(t time.Time) string {
				if t.IsZero() {
					return "never"
				}
				if dur := time.Since(t); dur < time.Duration(5*24)*time.Hour {
					return timediff.TimeDiff(t)
				}
				return t.Format("2006-01-02 15:04:05 MST")
			},
			"formatDurationSec": func(durationInSecs int) string {
				hrs := durationInSecs / 3600
				mins := (durationInSecs / 60) % 60
				secs := durationInSecs % 60
				if hrs > 0 {
					return fmt.Sprintf("%d:%02d:%02d", hrs, mins, secs)
				}
				return fmt.Sprintf("%d:%02d", mins, secs)
			},
		}),
	).Use(handler)

	closeFn = func() {
		jobStore.Close()
		videoStore.Close()
		feedItemStore.Close()
		feedsStore.Close()
	}

	return handler, closeFn, nil
}
