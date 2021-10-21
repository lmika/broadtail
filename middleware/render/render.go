package render

import (
	"bytes"
	"net/http"
)

func HTML(r *http.Request, w http.ResponseWriter, status int, templateName string) {
	rc, ok := r.Context().Value(renderContextKey).(*renderContext)
	if !ok {
		return
	}

	bw := new(bytes.Buffer)
	if err := rc.config.Templates.ExecuteTemplate(bw, templateName, rc.values); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	bw.WriteTo(w)
}
