.PHONY: build test mock clean

APP_NAME=stress-test

build:
	go build -o bin/$(APP_NAME) .

test:
	go test ./...

mock:
	mockery

clean:
	rm -rf bin/$(APP_NAME)
