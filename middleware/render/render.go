package render

import (
	"bytes"
	"html/template"
	"net/http"
)

func HTML(r *http.Request, w http.ResponseWriter, status int, templateName string) {
	rc, ok := r.Context().Value(renderContextKey).(*renderContext)
	if !ok {
		return
	}

	// Render the content template
	tmpl, err := rc.config.template(templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bw := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(bw, templateName, rc.values); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the master template
	masterTmpl, err := rc.config.template("masters/frame.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	masterBw := new(bytes.Buffer)
	if err := masterTmpl.ExecuteTemplate(masterBw, "masters/frame.html", struct{
		Content template.HTML
	}{
		Content: template.HTML(bw.String()),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	masterBw.WriteTo(w)
}
