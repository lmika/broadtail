module github.com/lmika/broadtail

go 1.18

require (
	github.com/asdine/storm/v3 v3.2.1
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/kkyr/fig v0.3.0
	github.com/lmika/gopkgs v0.0.0-20220504060120-48c3e4f681e8
	github.com/lmika/shellwords v0.0.0-20140714114018-ce258dd729fe
	github.com/mergestat/timediff v0.0.2
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/net v0.0.0-20211029224645-99673261e6eb // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/lmika/gopkgs => ../../libs/gopkgs
