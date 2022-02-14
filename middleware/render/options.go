package render

import "html/template"

type ConfigOption func(config *Config)

func WithFuncs(funcs template.FuncMap) ConfigOption {
	return func(config *Config) {
		config.funcMaps = funcs
	}
}
