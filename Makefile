.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/getDeviceBin src/getDevice/getDevice.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/createDeviceBin src/createDevice/createDevice.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

test:
	go test ./src/getDevice/
	go test ./src/createDevice/
