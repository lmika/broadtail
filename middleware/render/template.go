package render

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"sync"
)

type Config struct {
	templateFS fs.FS
	useCache bool

	cacheMutex *sync.RWMutex
	templateCache map[string]*template.Template
}

func New(tmplFS fs.FS, useCache bool) *Config {
	return &Config{
		templateFS: tmplFS,
		useCache: useCache,

		cacheMutex: new(sync.RWMutex),
		templateCache: make(map[string]*template.Template),
	}
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

func (tc *Config) template(name string) (*template.Template, error) {
	if !tc.useCache {
		return tc.parseTemplate(name)
	}

	tc.cacheMutex.RLock()
	tmpl, hasTmpl := tc.templateCache[name]
	tc.cacheMutex.RUnlock()

	if hasTmpl {
		return tmpl, nil
	}
	parsedTmpl, err := tc.parseTemplate(name)
	if err != nil {
		return nil, err
	}

	tc.cacheMutex.Lock()
	tc.templateCache[name] = parsedTmpl
	tc.cacheMutex.Unlock()

	return parsedTmpl, nil
}

func (tc *Config) parseTemplate(name string) (*template.Template, error) {
	f, err := tc.templateFS.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tmplBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New(name).Parse(string(tmplBytes))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
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
