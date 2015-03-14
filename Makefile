ifdef DATA_DIR
	OPTIONS=--data-dir $(DATA_DIR)
endif


all: run

clean:
	rm -f resources.go wiki

fmt:
	go fmt

get:
	go get

resources.go: resources/*
	$(GOPATH)/bin/esc -o resources.go -prefix="resources" `find resources \( -name "*.css" -o -name "*.js" -o -name "*.html" \)`
	go fmt resources.go

run: resources.go
	go run git_repo.go git_storage.go handlers.go templates.go resources.go wiki.go $(OPTIONS)

wiki: *.go
	go build


.PHONY: all clean fmt get run
