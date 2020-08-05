EXECUTABLES := curl go golangci-lint
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))


.PHONY: lint
lint:
	sed -i 's/[ \t]*$$//' $(shell find . -name "*.md")
	sed -i 's/[ \t]*$$//' $(shell find . -name "*.go")
	golangci-lint run

cover: test codecov_upload.sh
ifeq ("$(wildcard .codecov_token)","")
	./codecov_upload.sh -c -t "$${CODECOV_TOKEN}" -C "$${COMMIT_SHA}" -Z
else
	./codecov_upload.sh -c -t "$(shell cat .codecov_token)" -Z
endif


.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out


.PHONY: build
build:  clean test cover lint
	go build -o app ./...


codecov_upload.sh:
	curl -s https://codecov.io/bash -o codecov_upload.sh
	chmod +x ./codecov_upload.sh
