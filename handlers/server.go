package handlers

import (
	"context"
	"fmt"
	"github.com/lmika/broadtail/services/feedfetchers/appledevvideos"
	"html"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/youtubedl"
	"github.com/lmika/broadtail/services/feedfetchers"
	"github.com/lmika/broadtail/services/feedfetchers/youtuberss"
	"github.com/lmika/broadtail/services/videosources/youtubevideosource"

	"github.com/lmika/broadtail/handlers/monitor"
	"github.com/lmika/broadtail/handlers/settings"
	"github.com/lmika/broadtail/services/videodownload"
	"github.com/lmika/broadtail/services/videosources"
	"github.com/lmika/broadtail/services/videosources/simulatorvideosource"

	"github.com/gorilla/websocket"
	"github.com/lmika/broadtail/providers/plexprovider"
	"github.com/lmika/broadtail/services/favourites"
	"github.com/lmika/broadtail/services/rules"
	"github.com/lmika/gopkgs/http/middleware/render"

	"github.com/lmika/broadtail/services/videomanager"
	"github.com/mergestat/timediff"
	"github.com/robfig/cron"

	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/jobdispatcher"
	"github.com/lmika/broadtail/middleware/ujs"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/lmika/broadtail/providers/stormstore"
	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/pkg/errors"
)

type Config struct {
	LibraryDir   string
	LibraryOwner string

	BaseDataDir        string
	JobDataFile        string
	VideoDataFile      string
	FeedsDataFile      string
	FeedItemsDataFile  string
	FavouritesDataFile string

	PlexBaseURL string
	PlexToken   string

	YTDownloadCommand   []string
	YTDownloadSimulator bool

	TemplateFS fs.FS
	AssetFS    fs.FS
}

func Server(config Config) (handler http.Handler, closeFn func(), err error) {
	dispatcher := jobs.New()
	dbManager := stormstore.NewDBManager(config.BaseDataDir)
	defer dbManager.Close()

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
	favouriteStore, err := stormstore.NewFavouriteStore(config.FavouritesDataFile)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open favourites store")
	}
	rulesStore, err := stormstore.NewRulesStore(dbManager)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot open rules store")
	}

	youtubeRssFetcher := youtuberss.New()
	feedFetcher := feedfetchers.NewService(map[string]feedfetchers.FeedDriver{
		models.FeedTypeYoutubeChannel:  youtubeRssFetcher,
		models.FeedTypeYoutubePlaylist: youtubeRssFetcher,
		models.FeedTypeAppleDev:        appledevvideos.NewService(),
	})

	var youtubedlProvider videosources.SourceProvider
	if config.YTDownloadSimulator {
		log.Println("Using youtuble-dl simulator")
		youtubedlProvider = simulatorvideosource.NewService()
	} else {
		yp, err := youtubedl.New(config.YTDownloadCommand)
		if err != nil {
			return nil, nil, errors.Wrap(err, "cannot instantiate youtube-dl provider")
		}
		youtubedlProvider = youtubevideosource.NewService(yp)
	}

	videoSourcesServices := videosources.NewService(map[models.VideoRefSource]videosources.SourceProvider{
		models.YoutubeVideoRefSource: youtubedlProvider,
	})

	plexProvider := plexprovider.New(config.PlexBaseURL, config.PlexToken)

	jobsManager := jobsmanager.New(dispatcher, jobStore)
	vidDownloadService := videodownload.NewService(videodownload.Config{
		LibraryDir:          config.LibraryDir,
		LibraryOwner:        config.LibraryOwner,
		VideoSourcesService: videoSourcesServices,
		VideoStore:          videoStore,
		VideoDownloadHooks:  plexProvider,
		FeedStore:           feedsStore,
		FeedItemStore:       feedItemStore,
		JobsManager:         jobsManager,
	})

	favouriteService := favourites.NewService(favouriteStore, vidDownloadService, feedsStore, feedItemStore)
	feedsManager := feedsmanager.New(feedsStore, feedItemStore, feedFetcher, favouriteService, rulesStore, vidDownloadService)
	videoManager := videomanager.New(config.LibraryDir, videoStore)
	rulesService := rules.NewService(rulesStore, feedsStore)

	go jobsManager.Start()

	// Schedule updates every hour
	c := cron.New()
	if err := c.AddFunc("@every 15m", func() {
		if err := feedsManager.UpdateAllFeeds(context.Background()); err != nil {
			log.Printf("error updating all feeds: %v", err)
		}
	}); err != nil {
		return nil, nil, errors.Wrap(err, "invalid feed update schedule")
	}
	c.Start()

	indexHandlers := &indexHandlers{jobsManager: jobsManager, feedsManager: feedsManager, upgrader: websocket.Upgrader{}}
	ytdownloadHandlers := &youTubeDownloadHandlers{videoDownloadService: vidDownloadService, jobsManager: jobsManager}
	detailsHandler := &detailsHandler{videoSources: videoSourcesServices, videoManager: videoManager, favouriteService: favouriteService}
	videoHandler := &monitor.VideoHandlers{VideoManager: videoManager}
	jobsHandlers := &monitor.JobsHandlers{JobsManager: jobsManager}
	feedsHandlers := &feedsHandler{feedsManager: feedsManager}
	favouritesHandlers := &favouritesHandler{favouriteService: favouriteService}

	settingIndexHandlers := &settings.IndexHandlers{}
	settingRulesHandlers := &settings.RulesHandlers{RulesService: rulesService, FeedManager: feedsManager}

	r := mux.NewRouter()
	r.Handle("/", indexHandlers.Index()).Methods("GET")
	r.Handle("/ws/status", indexHandlers.StatusUpdateWebsocket()).Methods("GET")
	r.Handle("/job/download/video", ytdownloadHandlers.CreateDownloadJob()).Methods("POST")

	r.Handle("/quicklook", detailsHandler.QuickLook()).Methods("GET")
	r.Handle("/details/video/{video_id}", detailsHandler.VideoDetails()).Methods("GET")

	r.Handle("/monitor/videos", videoHandler.List()).Methods("GET")
	r.Handle("/monitor/videos/{video_id}", videoHandler.Show()).Methods("GET")

	r.Handle("/monitor/jobs", jobsHandlers.List()).Methods("GET")
	r.Handle("/monitor/jobs/done", jobsHandlers.ClearDone()).Methods("DELETE")
	r.Handle("/monitor/jobs/{job_id}", jobsHandlers.Show()).Methods("GET")
	r.Handle("/monitor/jobs/{job_id}", jobsHandlers.Delete()).Methods("DELETE")

	feedsHandlers.Routes(r)
	favouritesHandlers.Routes(r)
	settingIndexHandlers.Routes(r)
	settingRulesHandlers.Routes(r)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.FS(config.AssetFS))))

	handler = ujs.MethodRewriteHandler(r)
	handler = jobdispatcher.New(dispatcher).Use(handler)
	handler = render.New(
		config.TemplateFS,
		render.WithFuncs(template.FuncMap{
			"mup": func(x, y float64) int {
				return int(x * y)
			},
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
			"selectOption": func(value, optionValue any, name string) template.HTML {
				escapedValue := html.EscapeString(fmt.Sprint(optionValue))
				escapedName := html.EscapeString(fmt.Sprint(name))
				if reflect.DeepEqual(value, optionValue) {
					return template.HTML(fmt.Sprintf(`<option value="%s" selected>%s</option>`, escapedValue, escapedName))
				}

				return template.HTML(fmt.Sprintf(`<option value="%s">%s</option>`, escapedValue, escapedName))
			},
			"checkbox": func(value bool, name string, label string) template.HTML {
				escapedName := html.EscapeString(name)
				escapedLabel := html.EscapeString(label)
				checkAttr := ""
				if value {
					checkAttr = "checked"
				}

				return template.HTML(fmt.Sprintf(`
					<input name="%s" type="hidden" value="off">
                	<label>
                    	<input name="%s" type="checkbox" value="on" %s> %s
                	</label>
				`, escapedName, escapedName, checkAttr, escapedLabel))
			},
		}),
		render.WithFrame("frames/main.html"),
		render.RebuildOnChange(context.Background(), "templates"), // TEMP
	).Use(handler)

	closeFn = func() {
		jobStore.Close()
		videoStore.Close()
		feedItemStore.Close()
		feedsStore.Close()
		favouriteStore.Close()
	}

	return handler, closeFn, nil
}
