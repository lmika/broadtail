package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/feedsmanager"
)

type feedsHandler struct {
	feedsManager *feedsmanager.FeedsManager
}

func (h *feedsHandler) List() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feeds, err := h.feedsManager.List(ctx)
		if err != nil {
			return err

		}
		render.Set(r, "feeds", feeds)
		render.HTML(r, w, http.StatusOK, "feeds/index.html")
		return nil
	})
}

func (h *feedsHandler) New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var feed models.Feed

		if err := reqbind.Bind(&feed, r); err != nil {
			feed = models.Feed{}
		}

		render.Set(r, "feed", feed)
		render.HTML(r, w, http.StatusOK, "feeds/new.html")
	})
}

func (h *feedsHandler) Create() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var feed models.Feed

		if err := reqbind.Bind(&feed, r); err != nil {
			return err
		}

		if err := h.feedsManager.Save(ctx, &feed); err != nil {
			return err
		}

		http.Redirect(w, r, "/feeds/"+feed.ID.String(), http.StatusSeeOther)
		return nil
	})
}

func (h *feedsHandler) Show() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["feed_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		feed, err := h.feedsManager.Get(ctx, feedId)
		if err != nil {
			return err
		}

		recentItems, err := h.feedsManager.RecentFeedItems(ctx, feedId)
		if err != nil {
			return err
		}

		render.Set(r, "feed", feed)
		render.Set(r, "recentItems", recentItems)
		render.HTML(r, w, http.StatusOK, "feeds/show.html")
		return nil
	})
}

func (h *feedsHandler) Refresh() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["feed_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		if err := h.feedsManager.UpdateFeed(ctx, feedId); err != nil {
			return err
		}

		http.Redirect(w, r, "/feeds/"+feedId.String(), http.StatusSeeOther)
		return nil
	})
}

func (h *feedsHandler) Edit() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["feed_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		feed, err := h.feedsManager.Get(ctx, feedId)
		if err != nil {
			return err
		}

		render.Set(r, "feed", feed)
		render.HTML(r, w, http.StatusOK, "feeds/edit.html")
		return nil
	})
}

func (h *feedsHandler) Update() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["feed_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		feed, err := h.feedsManager.Get(ctx, feedId)
		if err != nil {
			return err
		}

		if err := reqbind.Bind(&feed, r); err != nil {
			return err
		}

		if err := h.feedsManager.Save(ctx, &feed); err != nil {
			return err
		}

		http.Redirect(w, r, "/feeds/"+feed.ID.String(), http.StatusSeeOther)
		return nil
	})
}

func (h *feedsHandler) Delete() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		feedId, err := uuid.Parse(mux.Vars(r)["feed_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed ID: %v", err.Error())
		}

		if err := h.feedsManager.Delete(ctx, feedId); err != nil {
			return err
		}

		http.Redirect(w, r, "/feeds", http.StatusSeeOther)
		return nil
	})
}

func (h *feedsHandler) Routes(r *mux.Router) {
	r.Handle("/feeds", h.List()).Methods("GET")
	r.Handle("/feeds/new", h.New()).Methods("GET")
	r.Handle("/feeds", h.Create()).Methods("POST")
	r.Handle("/feeds/{feed_id}", h.Show()).Methods("GET")
	r.Handle("/feeds/{feed_id}/refresh", h.Refresh()).Methods("POST")
	r.Handle("/feeds/{feed_id}/edit", h.Edit()).Methods("GET")
	r.Handle("/feeds/{feed_id}", h.Update()).Methods("PUT", "POST")
	r.Handle("/feeds/{feed_id}", h.Delete()).Methods("DELETE")
}
