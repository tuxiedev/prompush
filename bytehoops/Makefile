.PHONY: help build test docker test

help:
	@echo "Build bytehoops project"
	@echo "	test			runs tests"
	@echo "	build			builds a bianry for bytehoops"
	@echo "	docker			builds a docker image for bytehoops"
	@echo "	clean			runs go clean"

build:
	go build

docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
	docker build -t tuxiedev/bytehoops .

test:
	docker-compose -f test-services.yml up -d
	go test ./...
	docker-compose -f test-services.yml down
clean:
	docker-compose -f test-services.yml down
	go clean
