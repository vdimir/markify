package app

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/vdimir/markify/app/apperr"
	"github.com/vdimir/markify/app/view"
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

	app.addFileServer(r, "/assets")
	app.addFixedPages(r)

	r.Get("/ping", app.handlePing)
	r.Get("/robots.txt", app.handleRobotsTxt)

	r.Get("/", app.handlePageIndex)

	r.Get("/p/{pageID}", app.handleViewPageDoc)

	r.Get("/link", app.handlePageInputURL)
	r.Post("/link", app.handleCreateDocument)

	r.Get("/compose", app.handlePageTextInput)
	r.Post("/compose", app.handleCreateDocument)

	r.Post("/preview", app.handlePagePreview)
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

func (app *App) handlePageInputURL(w http.ResponseWriter, r *http.Request) {
	ctx := &view.URLPromptContext{
		Title: defaultTitle,
	}
	app.viewTemplate(http.StatusOK, ctx, w)
}

func (app *App) handlePageTextInput(w http.ResponseWriter, r *http.Request) {
	ctx := &view.EditorContext{
		Title: defaultTitle,
	}
	app.viewTemplate(http.StatusOK, ctx, w)
}

func (app *App) handleCreateDocument(w http.ResponseWriter, r *http.Request) {
	createReq := parseUserInput(r)
	docID, err := app.saveDocument(createReq)
	if err != nil {
		if errUser, ok := err.(apperr.UserError); ok {
			var returnToPageCtx view.TemplateContext
			if createReq.IsURL {
				returnToPageCtx = &view.URLPromptContext{
					Title: fmt.Sprintf("%s :(", defaultTitle),
					Msg:   errUser.String(),
				}
			} else {
				returnToPageCtx = &view.EditorContext{
					Title:       fmt.Sprintf("%s :(", defaultTitle),
					Msg:         errUser.String(),
					InitialText: string(createReq.Data),
				}
			}
			app.viewTemplate(http.StatusBadRequest, returnToPageCtx, w)
		} else {
			app.serverError(err, w)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/p/%s", docID), 302)
}

func (app *App) handleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	allowedPaths := []string{
		"/$", "/help$",
	}
	buf := bytes.NewBufferString("User-agent: *\nDisallow: /\n")
	for _, path := range allowedPaths {
		buf.WriteString(fmt.Sprintf("Allow: %s\n", path))
	}
	buf.WriteTo(w)
}

func (app *App) handlePagePreview(w http.ResponseWriter, r *http.Request) {
	createReq := parseUserInput(r)
	doc, err := app.engine.CreateDocument(createReq)
	if err != nil {
		app.serverError(err, w)
		return
	}
	app.viewDocument(doc, "Preview", w)
}

func (app *App) handleViewPageDoc(w http.ResponseWriter, r *http.Request) {
	pageID := chi.URLParam(r, "pageID")
	doc, err := app.getDocument([]byte(pageID))
	if err != nil {
		panic(err) // TODO
	}
	if doc == nil {
		app.notFound(w, r)
		return
	}
	app.viewDocument(doc, "", w)
}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	ctx := &view.StatusContext{
		Title:     "Not Found",
		HeaderMsg: "404",
		Msg:       "Page Not Found",
	}
	app.viewTemplate(http.StatusNotFound, ctx, w)
}
