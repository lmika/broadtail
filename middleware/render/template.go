package render

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Config struct {
	templateFS fs.FS

	cacheMutex  *sync.RWMutex
	templateSet *template.Template
	funcMaps    template.FuncMap
}

func New(tmplFS fs.FS, opts ...ConfigOption) *Config {
	cfg := &Config{
		templateFS: tmplFS,

		cacheMutex:  new(sync.RWMutex),
		templateSet: template.New("/"),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.funcMaps != nil {
		cfg.templateSet = cfg.templateSet.Funcs(cfg.funcMaps)
	}

	_ = fs.WalkDir(tmplFS, ".", func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		tmpl, err := cfg.parseTemplate(path)
		if err != nil {
			log.Printf("template %v: %v", path, err)
			return nil
		}

		if _, err := cfg.templateSet.AddParseTree(path, tmpl.Tree); err != nil {
			log.Printf("template %v: %v", path, err)
		}

		return nil
	})
	return cfg
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
	return tc.templateSet.Lookup(name), nil
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

	tmpl := template.New(name)
	if tc.funcMaps != nil {
		tmpl = tmpl.Funcs(tc.funcMaps)
	}

	tmpl, err = tmpl.Parse(string(tmplBytes))
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
