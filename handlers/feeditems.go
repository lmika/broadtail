package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/render"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/services/feedsmanager"
)

type feedItemsHandler struct {
	feedsManager *feedsmanager.FeedsManager
}

func (h *feedItemsHandler) Update() http.Handler {
	type feedItemPatchReq struct {
		Favourite bool `json:"favourite"`
	}

	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var patchReq feedItemPatchReq

		feedId, err := uuid.Parse(mux.Vars(r)["item_id"])
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid feed item ID: %v", err.Error())
		}

		if err := reqbind.Bind(&patchReq, r); err != nil {
			return err
		}

		feedItem, err := h.feedsManager.GetFeedItem(ctx, feedId)
		if err != nil {
			return err
		} else if feedItem == nil {
			return errhandler.Errorf(http.StatusNotFound, "feed item not found")
		}

		feedItem.Favourite = patchReq.Favourite

		if err := h.feedsManager.SaveFeedItem(ctx, feedItem); err != nil {
			return err
		}

		render.JSON(r, w, http.StatusOK, feedItem)
		return nil
	})
}

func (h *feedItemsHandler) Routes(r *mux.Router) {
	r.Handle("/feeditems/{item_id}", h.Update()).Methods("PATCH")
}
