.PHONY:	all clean code-vet code-fmt test get

DEPS := $(shell find . -type f -name "*.go" -printf "%p ")

all: code-vet code-fmt test build/argocd-vault-replacer

clean:
	$(RM) -rf build

get: $(DEPS)
	go get ./...

test: get
	go test ./...

build/argocd-vault-replacer: $(DEPS) get
	mkdir -p build
	CGO_ENABLED=0 go build -o build ./...

code-vet: $(DEPS) get
## Run go vet for this project. More info: https://golang.org/cmd/vet/
	@echo go vet
	go vet $$(go list ./... )

code-fmt: $(DEPS) get
## Run go fmt for this project
	@echo go fmt
	go fmt $$(go list ./... )

lint: $(DEPS) get
## Run golint for this project
	@echo golint
	golint $$(go list ./... )
