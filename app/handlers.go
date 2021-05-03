package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	r.Route("/api", func(rapi chi.Router) {
		rapi.Post("/preview", app.handleApiPreview)
		rapi.Post("/create", app.handleApiCreateDocument)
		rapi.NotFound(app.handleApiNotFound)
	})

	r.Get("/_ping", app.handlePing)
	r.Get("/ping", app.handlePing)
	r.Get("/_admin/unload", app.handleUnload)

	r.Get("/robots.txt", app.handleRobotsTxt)

	r.Get("/", app.handlePageIndex)

	r.Get("/p/{pageID}", app.handleViewPageDoc)
	r.Get("/p/{pageID}/text", app.handleViewPlainText)

	r.Get("/create", app.handlePageTextInput)
	r.Post("/create", app.handleCreateDocument)

	r.Post("/preview", app.handlePagePreview)
	r.Get("/preview", app.notFound)
	r.NotFound(app.notFound)
	return r
}

func (app *App) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if app.cfg.StatusText == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"status":"unloaded"}`))
		return
	}

	w.Write([]byte(app.cfg.StatusText))
}

func (app *App) handleUnload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	authHeader := r.Header.Get("Authorization")

	validPass := app.cfg.AdminPassword != "" &&
		strings.HasPrefix(authHeader, "Basic ") &&
		strings.TrimPrefix(authHeader, "Basic ") == app.cfg.AdminPassword

	if !validPass {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	app.cfg.StatusText = ""
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
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
	allowedPaths := []string{"/$", "/about$", "/info/*"}
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
	doc, err := app.converter.Convert(strings.NewReader(createReq.Text), createReq.Syntax)
	if err != nil {
		app.serverError(err, w)
		return
	}
	app.viewDocument(&Document{*doc, "", time.Time{}}, "Preview", "", w)
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

func (app *App) handleViewPlainText(w http.ResponseWriter, r *http.Request) {
	pageID := chi.URLParam(r, "pageID")
	data, _, err := app.blobStore.GetBlob(pageID)
	if err != nil {
		app.serverError(err, w)
		return
	}
	if data == nil {
		app.notFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = io.Copy(w, data)
	if err != nil {
		log.Printf("[ERROR] can't write response: %s", err.Error())
		return
	}
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
	req, err := parseCreateRequestForm(r)
	if err != nil {
		return nil, WrapfUserError(err, err.Error())
	}
	err = app.validatePasteRequest(req)
	if err != nil {
		return nil, WrapfUserError(err, err.Error())
	}
	return req, nil
}
