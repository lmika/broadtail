package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lmika/broadtail/middleware/errhandler"
	"github.com/lmika/broadtail/middleware/reqbind"
	"github.com/lmika/broadtail/models"
	"github.com/lmika/broadtail/services/favourites"
	"github.com/lmika/gopkgs/http/middleware/render"
	"github.com/pkg/errors"
)

type favouritesHandler struct {
	favouriteService *favourites.Service
}

func (fh *favouritesHandler) list() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var request = struct {
			Query string `req:"q"`
			Page  int    `req:"page"`
		}{}

		if err := reqbind.Bind(&request, r); err != nil {
			return err
		}

		feedItemFilter := models.ParseFeedItemFilter(request.Query)

		favouriteItems, err := fh.favouriteService.List(ctx, feedItemFilter, request.Page)
		if err != nil {
			return err
		}

		render.Set(r, "request", request)
		render.Set(r, "favourites", favouriteItems)
		render.HTML(r, w, http.StatusOK, "favourites/list.html")
		return nil
	})
}

func (fh *favouritesHandler) add() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		var req struct {
			Origin models.FavouriteOrigin `json:"origin"`
		}

		if err := reqbind.Bind(&req, r); err != nil {
			return err
		}

		log.Printf("Favourte request = %#v", req)

		favourite, err := fh.favouriteService.FavoriteVideoByOrigin(ctx, req.Origin)
		if err != nil {
			return errors.Wrap(err, "cannot add favourite")
		}

		render.JSON(r, w, http.StatusOK, favourite)
		return nil
	})
}

func (fh *favouritesHandler) delete() http.Handler {
	return errhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		favouriteIDStr, ok := mux.Vars(r)["favourite_id"]
		if !ok {
			return errhandler.Errorf(http.StatusBadRequest, "invalid favourite ID: %v", favouriteIDStr)
		}

		favouriteID, err := uuid.Parse(favouriteIDStr)
		if err != nil {
			return errhandler.Errorf(http.StatusBadRequest, "invalid favourite ID: %v", favouriteIDStr)
		}

		if err := fh.favouriteService.DeleteFavourite(ctx, favouriteID); err != nil {
			return errors.Wrapf(err, "cannot delete favourite ID: %v", favouriteID)
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	})
}

func (fh *favouritesHandler) Routes(r *mux.Router) {
	r.Handle("/favourites", fh.list()).Methods("GET")
	r.Handle("/favourites", fh.add()).Methods("POST")
	r.Handle("/favourites/{favourite_id}", fh.delete()).Methods("DELETE")
}
