export PROJECT_ROOT=$(shell pwd)

shortener-run:
	@go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/server/main.go