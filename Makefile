DEPENDENCIES = \
  github.com/gorilla/mux \
  github.com/mjibson/esc \
  github.com/russross/blackfriday

ifdef DATA_DIR
	OPTIONS=--data-dir $(DATA_DIR)
endif


all: run

bootstrap:
	go get $(DEPENDENCIES)

clean:
	rm -f resources.go wiki

fmt:
	go fmt

resources.go: resources/*
	$(GOPATH)/bin/esc -o resources.go -prefix="resources" `find resources \( -name "*.css" -o -name "*.js" -o -name "*.html" \)`

run: resources.go
	go run git_repo.go git_storage.go handlers.go templates.go resources.go wiki.go $(OPTIONS)

wiki: *.go
	go build


.PHONY: all bootstrap clean fmt run
