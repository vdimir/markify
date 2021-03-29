package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/vdimir/markify/view"
)

const defaultTitle = "markify"

// Routes setup
func (app *App) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)

	logFmt := &middleware.DefaultLogFormatter{
		Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: true}
	loggerMiddleware := middleware.RequestLogger(logFmt)
	r.Use(loggerMiddleware)

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)

	app.addFileServer(r, "public")
	app.addFixedPages(r)

	r.Get("/ping", app.handlePing)
	r.Get("/robots.txt", app.handleRobotsTxt)

	r.Get("/", app.handlePageIndex)

	r.Get("/p/{pageID}", app.handleViewPageDoc)

	r.Get("/create", app.handlePageTextInput)
	r.Post("/create", app.handleCreateDocument)

	r.Post("/preview", app.handlePagePreview)
	r.Get("/preview", app.notFound)
	r.NotFound(app.notFound)
	return r
}

func (app *App) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(app.cfg.StatusText))
}

func (app *App) handlePageIndex(w http.ResponseWriter, r *http.Request) {
	app.handlePageTextInput(w, r)
}


func (app *App) handlePageTextInput(w http.ResponseWriter, r *http.Request) {
	ctx := &view.EditorContext{
		Title: defaultTitle,
	}
	app.viewTemplate(http.StatusOK, ctx, w)
}

func (app *App) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	req, err := app.parseAndValidateRequest(r)
	if err != nil {
		app.respondError(err, req, w)
		return
	}
	docID, err := app.savePaste(req)
	if err != nil {
		app.respondError(err, req, w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/p/%s", docID), 302)
}

func (app *App) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	allowedPaths := []string{
		"/$", "/help$",
		"/p/*",
	}
	buf := bytes.NewBufferString("User-agent: *\nDisallow: /\n")
	for _, path := range allowedPaths {
		buf.WriteString(fmt.Sprintf("Allow: %s\n", path))
	}
	buf.WriteTo(w)
}

func (app *App) handlePagePreview(w http.ResponseWriter, r *http.Request) {
	createReq, err := app.parseAndValidateRequest(r)
	if err != nil {
		app.respondError(err, createReq, w)
		return
	}
	if err = app.converter.SupportSyntax(createReq.Syntax); err != nil {
		ctx := &view.StatusContext{
			Title:     "markify",
			HeaderMsg: "markify",
			Msg:       err.Error(),
		}
		app.viewTemplate(http.StatusBadRequest, ctx, w)
	}

	doc, err := app.converter.Convert(strings.NewReader(createReq.Text), createReq.Syntax)
	if err != nil {
		app.serverError(err, w)
		return
	}
	app.viewDocument(doc, "Preview", "", w)
}

func (app *App) handleViewPageDoc(w http.ResponseWriter, r *http.Request) {
	pageID := chi.URLParam(r, "pageID")
	doc, err := app.getDocument(pageID)
	if err != nil {
		app.serverError(err, w)
		return
	}
	if doc == nil {
		app.notFound(w, r)
		return
	}
	app.viewDocument(doc, "", r.URL.Path, w)
}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	ctx := &view.StatusContext{
		Title:     "Not Found",
		HeaderMsg: "404",
		Msg:       "Page Not Found",
	}
	app.viewTemplate(http.StatusNotFound, ctx, w)
}

func (app *App) parseAndValidateRequest(r *http.Request) (*CreatePasteRequest, error) {
	req, err := ParseCreatePasteRequest(r)
	if err != nil {
		return nil, WrapfUserError(err, err.Error())
	}
	err = app.validatePasteRequest(req)
	if err != nil {
		return nil, WrapfUserError(err, err.Error())
	}
	return req, nil
}
