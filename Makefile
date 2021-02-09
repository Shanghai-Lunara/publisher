.PHONY: build mod proto dc

HARBOR_DOMAIN := $(shell echo ${HARBOR})
PROJECT := lunara-common
SERVER_IMAGE := "$(HARBOR_DOMAIN)/$(PROJECT)/publisher:latest"

build:
	-i docker image rm $(SERVER_IMAGE)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o publisher cmd/v1/scheduler/main.go
	cp cmd/v1/scheduler/Dockerfile . && docker build -t $(SERVER_IMAGE) .
	rm -f Dockerfile && rm -f publisher
	docker push $(SERVER_IMAGE)

mod:
	go mod download
	go mod tidy

proto:
	cd scripts && bash ./gen.sh api

dc:
	go mod vendor
	cd scripts && bash ./gen.sh deepcopy
	rm -rf vendor

run:
	go run ./cmd/v1/scheduler/main.go -v=4 -configPath=./cmd/v1/scheduler/fake.yaml