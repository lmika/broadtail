.Phony: clean build run

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
	go run . -dev -data ./build/testdata -ytdl-simulator

build-js: prep
	npm run build-js
	npm run build-css

prep:
	mkdir -p build/assets/js
	mkdir -p build/assets/css
	mkdir -p build/testdata