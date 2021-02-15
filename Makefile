.PHONY:	all clean code-vet code-fmt test

all: code-vet code-fmt test build/vault-replacer

clean:
	$(RM) -rf build

test:
	go test ./...

build/vault-replacer: $(PACKAGE)
	go build -o ./...

code-vet: $(GENERATED_FILES)
## Run go vet for this project. More info: https://golang.org/cmd/vet/
	@echo go vet
	go vet $$(go list ./... )

code-fmt: $(GENERATED_FILES)
## Run go fmt for this project
	@echo go fmt
	go fmt $$(go list ./... )
