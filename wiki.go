package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	pagesDir      = "pages"
	pageExtension = ".md"
)

func bodyAction(action string) func(r *http.Request, rm *mux.RouteMatch) bool {
	return func(r *http.Request, rm *mux.RouteMatch) bool {
		return r.FormValue("action") == action
	}
}

func pageNameMatcher(r *http.Request, rm *mux.RouteMatch) bool {
	return r.URL.Path != "/_" && !strings.HasPrefix(r.URL.Path, "/_/")
}

func main() {
	var addr string
	var dataDir string
	flag.StringVar(&addr, "addr", "127.0.0.1:8000", "TCP address to listen on")
	flag.StringVar(&dataDir, "data-dir", "data", "Data directory")
	flag.Parse()

	storage := NewGitStorage(dataDir, pagesDir, pageExtension)
	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	templates, err := NewTemplates()
	if err != nil {
		log.Fatal(err)
	}
	app := AppContext{
		Storage:   storage,
		templates: templates,
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	fileServer := http.FileServer(FS(false))
	router.PathPrefix("/_/static/").Handler(http.StripPrefix("/_/", fileServer))

	for _, path := range []string{"/", "/{title:.{1,}}"} {
		// Delete
		router.HandleFunc(path, app.deleteHandler).MatcherFunc(pageNameMatcher).Queries("action", "delete").Methods("POST")

		// History
		router.HandleFunc(path, app.historyHandler).MatcherFunc(pageNameMatcher).Queries("action", "history").Methods("GET")

		// Diff
		router.HandleFunc(path, app.diffHandler).
			MatcherFunc(pageNameMatcher).
			Methods("POST").
			Queries("action", "edit").
			MatcherFunc(bodyAction("diff"))

		// Preview
		router.HandleFunc(path, app.previewHandler).
			MatcherFunc(pageNameMatcher).
			Methods("POST").
			Queries("action", "edit").
			MatcherFunc(bodyAction("preview"))

		// Save
		router.HandleFunc(path, app.saveHandler).
			MatcherFunc(pageNameMatcher).
			Methods("POST").
			Queries("action", "edit").
			MatcherFunc(bodyAction("save"))

		// Edit
		router.HandleFunc(path, app.editHandler).MatcherFunc(pageNameMatcher).Queries("action", "edit").Methods("GET", "POST")

		// View
		router.HandleFunc(path, app.viewHandler).MatcherFunc(pageNameMatcher).Methods("GET")
	}

	router.HandleFunc("/_/deleted", app.deletedHandler).Methods("GET")
	router.HandleFunc("/_/pages", app.allPagesHandler).Methods("GET")
	router.HandleFunc("/_/search", app.searchHandler).Methods("GET")

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
