CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

swag:
	swag init -g ./api/router.go -o api/docs

tidy:
	go mod tidy
	go mod vendor