.Phony: init clean build run

init:
	go get
	npm install

clean:
	-go clean
	-rm -r build

build: build-js
	go build

build-linux: build-js
	GOOS=linux go build

run: prep
	(sleep 99999 | npm run watch-js &)
	(sleep 99999 | npm run watch-css &)
	go run . -dev -config ./build/config.yaml

run-sim: prep
	(sleep 99999 | npm run watch-js &)
	(sleep 99999 | npm run watch-css &)
	go run . -dev -config ./build/config.yaml -ytdl-simulator

build-js: prep
	npm run build-js
	npm run build-css

prep:
	mkdir -p build/assets/js
	mkdir -p build/assets/css
	mkdir -p build/testdata
	echo "data_dir: `pwd`/build/testdata" > build/config.yaml