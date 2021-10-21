package render

import (
	"context"
	"html/template"
	"net/http"
)

type Config struct {
	Templates *template.Template
}

func New(tmpl *template.Template) *Config {
	return &Config{tmpl}
}

func (tc *Config) Use(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &renderContext{
			config: tc,
			values: make(map[string]interface{}),
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), renderContextKey, rc)))
	})
}

func Set(r *http.Request, name string, value interface{}) {
	rc, ok := r.Context().Value(renderContextKey).(*renderContext)
	if !ok {
		return
	}
	rc.values[name] = value
}

type renderContext struct {
	config *Config
	values map[string]interface{}
}

type renderContextKeyType struct{}

var renderContextKey = renderContextKeyType{}
