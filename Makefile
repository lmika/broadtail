.Phony: init clean build run

init:
	go get
	npm install
	go install github.com/lmika/rwt@latest

clean:
	-go clean
	-rm -r build

build: build-js
	go build

build-linux: build-js
	GOOS=linux go build

run: prep
	rwt watch &
	go run . -dev -config ./build/config.yaml

run-sim: prep
	rwt watch &
	go run . -dev -config ./build/config.yaml -ytdl-simulator

build-js: prep
	rwt build

release: clean
	goreleaser release --skip-validate --skip-publish

prep:
	mkdir -p build/assets/js
	mkdir -p build/assets/css
	mkdir -p build/testdata
	echo "data_dir: `pwd`/build/testdata" > build/config.yaml
	echo "library_dir: `pwd`/build/library" >> build/config.yaml