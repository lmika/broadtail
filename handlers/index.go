package handlers

import (
	"log"
	"net/http"

	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/lmika/broadtail/services/jobsmanager"
)

type indexHandlers struct {
	jobsManager  *jobsmanager.JobsManager
	feedsManager *feedsmanager.FeedsManager
}

func (ih *indexHandlers) Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recentFeedItems, err := ih.feedsManager.RecentFeedItemsFromAllFeeds(r.Context())
		if err != nil {
			log.Printf("warn: cannot get list of recent feed items: %v", err)
		}

		render.Set(r, "recentFeedItems", recentFeedItems)
		render.Set(r, "jobs", ih.jobsManager.RecentJobs())
		render.HTML(r, w, http.StatusOK, "index.html")
	})
}
