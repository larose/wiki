package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

type AllPagesContext struct {
	PageContext
	Titles []string
	Error  string
}

type AppContext struct {
	Storage   *GitStorage
	templates map[string]*template.Template
}

type Commit struct {
	ID      string
	Date    string
	Message string
	Delete  bool
}

type DeletedContext struct {
	PageContext
	Titles []string
	Error  string
}

type EditContext struct {
	PageContext
	Body          template.HTML
	BodySource    string
	CommitMessage string
	Edit          bool
	Preview       bool
	Diff          bool
}

type HistoryContext struct {
	PageContext
	Commits []Commit
	Error   string
}

type Page struct {
	Title string
	Body  template.HTML
}

type PageContext struct {
	Title    string
	SubTitle string
}

type PageSearchResult struct {
	Title string
	Lines []string
}

type PrintableContext struct {
	Title string
	Body  template.HTML
}

type SearchContext struct {
	PageContext
	SearchResults []PageSearchResult
	Query         string
	Error         string
}

type ViewContext struct {
	PageContext
	Body     template.HTML
	Revision string
	RawBody  string
}

func (app AppContext) allPagesHandler(w http.ResponseWriter, r *http.Request) {
	titles, err := app.Storage.ListPages()

	if err != nil {
		ctx := AllPagesContext{
			PageContext: PageContext{
				Title: "All Pages",
			},
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		app.templates["allPages"].Execute(w, ctx)
		return
	}

	ctx := AllPagesContext{
		PageContext: PageContext{
			Title: "All Pages",
		},
		Titles: titles,
	}
	if err := app.templates["allPages"].Execute(w, ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app AppContext) deleteHandler(w http.ResponseWriter, r *http.Request) {
	var title = normalizePath(mux.Vars(r)["title"])

	if err := app.Storage.DeletePage(title); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app AppContext) deletedHandler(w http.ResponseWriter, r *http.Request) {
	var titles []string
	var err error
	if titles, err = app.Storage.ListDeletedPages(); err != nil {
		ctx := DeletedContext{
			PageContext: PageContext{
				Title: "Deleted Pages",
			},
			Error: err.Error(),
		}
		renderError(app.templates["deleted"], w, ctx, http.StatusInternalServerError)
		return
	}

	ctx := DeletedContext{
		PageContext: PageContext{
			Title: "Deleted Pages",
		},
		Titles: titles,
	}
	renderTemplate(app.templates["deleted"], w, ctx)
}

func (app AppContext) diffHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	message := r.FormValue("message")
	title := normalizePath(mux.Vars(r)["title"])

	var out []byte
	var err error
	out, err = app.Storage.Diff(title, body)

	if err != nil {
		log.Println(err)
		return
	}

	ctx := EditContext{
		PageContext: PageContext{
			Title:    title,
			SubTitle: "diff",
		},
		Body:          template.HTML(string(out)),
		BodySource:    body,
		CommitMessage: message,
		Diff:          true,
	}
	app.templates["diff"].Execute(w, ctx)
}

func (app AppContext) editHandler(w http.ResponseWriter, r *http.Request) {
	var title = normalizePath(mux.Vars(r)["title"])
	var ctx EditContext

	if r.Method == "GET" {
		body, err := app.Storage.PageBody(title, "HEAD")
		if err != nil {
			ctx = EditContext{
				PageContext: PageContext{
					Title:    title,
					SubTitle: "edit",
				},
				Edit: true,
			}
		} else {
			ctx = EditContext{
				PageContext: PageContext{
					Title:    title,
					SubTitle: "edit",
				},
				BodySource: string(body),
				Edit:       true,
			}
		}
	} else {
		body := r.FormValue("body")
		message := r.FormValue("message")
		ctx = EditContext{
			PageContext: PageContext{
				Title:    title,
				SubTitle: "edit",
			},
			BodySource:    string(body),
			CommitMessage: message,
			Edit:          true,
		}
	}

	app.templates["edit"].Execute(w, ctx)
}

func (app AppContext) historyHandler(w http.ResponseWriter, r *http.Request) {
	title := normalizePath(mux.Vars(r)["title"])

	commits, err := app.Storage.History(title)

	if err != nil {
		ctx := HistoryContext{
			PageContext: PageContext{
				Title:    title,
				SubTitle: "history",
			},
			Error: err.Error(),
		}
		renderError(app.templates["history"], w, ctx, http.StatusInternalServerError)
		return
	}

	ctx := HistoryContext{
		PageContext: PageContext{
			Title:    title,
			SubTitle: "history",
		},
		Commits: commits,
	}
	renderTemplate(app.templates["history"], w, ctx)
}

func normalizePath(path string) string {
	if path == "" {
		return "home"
	}

	return path
}

func (app AppContext) previewHandler(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	message := r.FormValue("message")
	title := normalizePath(mux.Vars(r)["title"])

	ctx := EditContext{
		PageContext: PageContext{
			Title:    title,
			SubTitle: "preview",
		},
		Body:          template.HTML(renderMarkdown([]byte(body))),
		BodySource:    body,
		CommitMessage: message,
		Preview:       true,
	}
	app.templates["preview"].Execute(w, ctx)
}

func renderError(t *template.Template, w http.ResponseWriter, ctx interface{}, s int) {
	w.WriteHeader(http.StatusInternalServerError)
	renderTemplate(t, w, ctx)
}

func renderMarkdown(source []byte) []byte {
	flags := 0
	flags |= blackfriday.HTML_TOC
	flags |= blackfriday.HTML_SAFELINK

	extensions := 0
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	extensions |= blackfriday.EXTENSION_AUTO_HEADER_IDS

	renderer := blackfriday.HtmlRenderer(flags, "", "")
	return blackfriday.Markdown(source, renderer, extensions)
}

func renderTemplate(t *template.Template, w http.ResponseWriter, ctx interface{}) {
	if err := t.Execute(w, ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app AppContext) searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	searchResults, err := app.Storage.Search(q)
	if err != nil {
		ctx := SearchContext{
			PageContext: PageContext{
				Title: "Search results for " + q,
			},
			Query: q,
			Error: err.Error(),
		}
		renderError(app.templates["search"], w, ctx, http.StatusInternalServerError)
		return
	}

	ctx := SearchContext{
		Query:         q,
		SearchResults: searchResults,
	}

	app.templates["search"].Execute(w, ctx)
}

func (app *AppContext) saveHandler(w http.ResponseWriter, r *http.Request) {
	var title = normalizePath(mux.Vars(r)["title"])

	body := r.FormValue("body")
	message := r.FormValue("message")

	if message == "" {
		message = "Update " + title
	}

	if err := app.Storage.SetPageBody(title, body, message); err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/"+title, http.StatusSeeOther)
}

func (app AppContext) viewHandler(w http.ResponseWriter, r *http.Request) {
	var title = normalizePath(mux.Vars(r)["title"])
	revision := r.URL.Query().Get("revision")
	format := r.URL.Query().Get("format")
	body, err := app.Storage.PageBody(title, revision)
	if err != nil {
		if dirty, ok := err.(*DirtyWorkTree); ok {
			w.Write([]byte(dirty.Error()))
			return
		}

		http.Redirect(w, r, "/"+title+"?action=edit", http.StatusFound)
		return
	}

	if format == "raw" {
		w.Write(body)
		return
	}

	if format == "printable" {
		ctx := PrintableContext{
			Title: title,
			Body:  template.HTML(renderMarkdown(body)),
		}
		app.templates["printable"].Execute(w, ctx)
		return
	}

	ctx := ViewContext{
		PageContext: PageContext{
			Title:    title,
			SubTitle: revision,
		},
		Body:     template.HTML(renderMarkdown(body)),
		Revision: revision,
		RawBody:  string(body),
	}

	app.templates["view"].Execute(w, ctx)
}
