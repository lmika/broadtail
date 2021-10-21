package handlers

import (
	"github.com/lmika/broadtail/middleware/render"
	"net/http"
)

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Set(r, "message", "Hello renderer")
		render.HTML(r, w, http.StatusOK, "index.html")
	})
}
