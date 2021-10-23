.Phony: clean build run

clean:
	-go clean
	-rm -r build

build: build-js
	go build

build-linux: build-js
	GOOS=linux go build

run: build-js
	go run . -dev

build-js: prep
	npm run build-js

prep:
	mkdir build/assets/js