package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/providers/jobs"
	"github.com/pkg/errors"

	"github.com/lmika/broadtail/services/feedsmanager"
	"github.com/lmika/broadtail/services/jobsmanager"
	"github.com/lmika/gopkgs/http/middleware/render"
)

type indexHandlers struct {
	jobsManager  *jobsmanager.JobsManager
	feedsManager *feedsmanager.FeedsManager
	upgrader     websocket.Upgrader
}

func (ih *indexHandlers) Index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recentFeedItems, err := ih.feedsManager.RecentFeedItemsFromAllFeeds(r.Context(), models.FeedItemFilter{}, 0, 10)
		if err != nil {
			log.Printf("warn: cannot get list of recent feed items: %v", err)
		}

		render.Set(r, "recentFeedItems", recentFeedItems)
		render.Set(r, "jobs", ih.jobsManager.RecentJobs())
		render.HTML(r, w, http.StatusOK, "index.html")
	})
}

func (ih *indexHandlers) StatusUpdateWebsocket() http.Handler {
	type wbJobUpdateMessage struct {
		ID      string  `json:"id"`
		Type    string  `json:"type"`
		State   string  `json:"state,omitempty"`
		Summary string  `json:"summary,omitempty"`
		Percent float64 `json:"percent"`
	}

	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		c, err := ih.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return errors.Wrap(err, "cannot update socket")
		}
		defer c.Close()

		sub := ih.jobsManager.Dispatcher().Subscribe()
		defer sub.Close()

		for msg := range sub.Chan() {
			var err error = nil

			switch m := msg.(type) {
			case jobs.UpdateSubscriptionEvent:
				err = c.WriteJSON(wbJobUpdateMessage{
					ID:      m.Job.ID().String(),
					Type:    "update",
					Summary: m.Update.Summary,
					Percent: m.Update.Percent,
				})
			case jobs.StateTransitionSubscriptionEvent:
				err = c.WriteJSON(wbJobUpdateMessage{
					ID:    m.Job.ID().String(),
					Type:  "newstate",
					State: m.ToState.String(),
				})
			}

			if err != nil {
				break
			}
		}
		return nil
	})
}
