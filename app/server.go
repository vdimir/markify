package app

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/vdimir/markify/render/htmltemplate"
	"github.com/vdimir/markify/render/md"
	"github.com/vdimir/markify/util"

	"github.com/pkg/errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const fixedPagesPrefixDir = "/static_pages"

// StartServer listen incoming requsets and hable it
func (app *App) StartServer() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	// r.Use(middleware.Recoverer)

	app.addFileServer(r, "/assets")
	app.addFixedPages(r)

	r.Get("/ping", app.handlePing)
	r.Get("/robots.txt", app.handleRobotsTxt)

	r.Get("/", app.handleIndexPage)
	r.Get("/p/{pageID}", app.handleReadHashPage)

	r.Get("/link", app.handleUrlInput)
	r.Post("/link", app.handleNewByURL)

	r.Get("/compose", app.handleCreateTextPage)
	r.Post("/compose", app.handleNewByText)
	r.Post("/preview", app.handlePreviewText)
	r.NotFound(app.notFound)

	serverURL := app.cfg.ServerAddrHost
	if serverURL == "" {
		serverURL = "localhost"
	}
	serverURL = fmt.Sprintf("http://%s:%d", serverURL, app.cfg.ServerPort)

	log.Printf("[INFO] starting server at %s\n", serverURL)

	app.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.cfg.ServerAddrHost, app.cfg.ServerPort),
		Handler: r,
	}

	err := app.httpServer.ListenAndServe()

	if err != nil {
		panic(errors.Wrap(err, "cannot start server"))
	}
}

func (app *App) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(app.cfg.StatusText))
}

func (app *App) handleIndexPage(w http.ResponseWriter, r *http.Request) {
	app.handleCreateTextPage(w, r)
}

func (app *App) handleUrlInput(w http.ResponseWriter, r *http.Request) {
	ctx := htmltemplate.InfoPageContext{
		Title: "markify",
	}
	app.wirteTemplateToResp(htmltemplate.URLInputPage, http.StatusOK, ctx, w)
}

func (app *App) handleCreateTextPage(w http.ResponseWriter, r *http.Request) {
	ctx := htmltemplate.InfoPageContext{
		Title: "markify",
	}
	app.wirteTemplateToResp(htmltemplate.TextEditorTemplate, http.StatusOK, ctx, w)
}

type formParams struct {
	TextData         []byte
	EnableShortcodes bool
	EnableRelImgLink bool
}

func parseFormParams(r *http.Request, dataKey string) formParams {
	return formParams{
		TextData:         []byte(r.FormValue(dataKey)),
		EnableShortcodes: r.FormValue("shortcode") == "on",
		EnableRelImgLink: r.FormValue("rel-link") == "on",
	}
}

func (app *App) handleNewByURL(w http.ResponseWriter, r *http.Request) {
	params := parseFormParams(r, "url")
	if len(params.TextData) == 0 {
		http.Redirect(w, r, "/compose", 302)
		return
	}

	urlHash, err := app.addPageByURL(params)
	if err != nil {
		if _, ok := err.(UserError); ok {
			ctx := htmltemplate.InfoPageContext{
				Title:   "markify",
				MainMsg: "Incorrect URL or data, try another:",
			}
			app.wirteTemplateToResp(htmltemplate.URLInputPage, http.StatusBadRequest, ctx, w)
		} else {
			app.serverError(err, w)
		}
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/p/%s", urlHash), 302)
}

func (app *App) handleNewByText(w http.ResponseWriter, r *http.Request) {
	params := parseFormParams(r, "text")
	urlHash, err := app.addPageByText(params)
	if err != nil {
		app.serverError(err, w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/p/%s", urlHash), 302)
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

func (app *App) handlePreviewText(w http.ResponseWriter, r *http.Request) {
	params := parseFormParams(r, "text")
	ropts := &md.Options{
		DisableShortcodes: !params.EnableShortcodes,
	}
	htmlBuf, err := app.render.Render(params.TextData, ropts)
	if err != nil {
		app.serverError(err, w)
	}

	ctx := htmltemplate.MdPageContext{
		Title: "Preview",
		Body:  template.HTML(htmlBuf.Bytes()),
	}
	app.wirteTemplateToResp(htmltemplate.ContentTemplatePage, http.StatusOK, ctx, w)
}

func (app *App) handleReadHashPage(w http.ResponseWriter, r *http.Request) {
	pageID := chi.URLParam(r, "pageID")
	app.handlePageFromCache([]byte(pageID), w)
}

func (app *App) handlePageFromCache(pageID []byte, w http.ResponseWriter) {
	data, err := app.pageCache.Load(pageID)

	if err != nil {
		app.serverError(err, w)
		return
	}
	if data == nil {
		app.notFound(w, nil)
		return
	}

	ctx := htmltemplate.MdPageContext{
		Title: string(pageID),
		Body:  template.HTML(data),
	}
	app.wirteTemplateToResp(htmltemplate.ContentTemplatePage, http.StatusOK, ctx, w)
}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	ctx := htmltemplate.InfoPageContext{
		Title:     "Not Found",
		HeaderMsg: "404",
		MainMsg:   "Page Not Found",
	}
	app.wirteTemplateToResp(htmltemplate.StatusTemplatePage, http.StatusNotFound, ctx, w)
}

func (app *App) serverError(err error, w http.ResponseWriter) {
	log.Printf("[ERROR] %v", err)
	ctx := htmltemplate.InfoPageContext{
		Title:     "Error",
		HeaderMsg: "500",
		MainMsg:   "Internal Server error",
	}
	app.wirteTemplateToResp(htmltemplate.StatusTemplatePage, http.StatusInternalServerError, ctx, w)
}

func (app *App) wirteTemplateToResp(path string, code int, pageContext interface{}, w http.ResponseWriter) {
	htmlBuf := &bytes.Buffer{}
	err := app.htmlTplRender.RenderPage(htmlBuf, path, pageContext)

	if err != nil {
		log.Printf("[ERROR] %v", errors.Wrapf(err, "cannot render template %s", path))
		app.serverErrorFallback(w)
		return
	}
	if code > 0 {
		w.WriteHeader(code)
	}
	_, err = htmlBuf.WriteTo(w)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return
	}
}

func (app *App) serverErrorFallback(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	rawHTML := "<!DOCTYPE html>" +
		"<html><head><title>Error</title></head><body>" +
		"Internal Server Error" +
		"</body></html>\n"
	w.Write([]byte(rawHTML))
}

func (app *App) addFixedPages(r chi.Router) {
	createDebugHanlder := func(filePath string) func(w http.ResponseWriter, r *http.Request) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			f, _ := app.staticFs.Open(filePath)
			data, _ := ioutil.ReadAll(f)
			htmlBuf, _ := app.render.Render(data, nil)
			_, name := path.Split(filePath)
			ctx := htmltemplate.MdPageContext{
				Title: name,
				Body:  template.HTML(htmlBuf.Bytes()),
			}
			app.wirteTemplateToResp(htmltemplate.ContentTemplatePage, http.StatusOK, ctx, w)
		}
		return handler
	}

	err := util.WalkFiles(app.staticFs, fixedPagesPrefixDir, func(data []byte, filePath string) error {
		name := strings.TrimSuffix(filePath, ".md")
		err := app.saveRenderToCache(data, []byte(name), nil)
		if err != nil {
			return errors.Wrapf(err, "add fixed pages: error rendering page")
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			app.handlePageFromCache([]byte(name), w)
		}

		rawHandler := func(w http.ResponseWriter, r *http.Request) {
			// w.WriteHeader(http.StatusOK)
			// w.Write(data)
			ctx := htmltemplate.MdPageContext{
				Title: name,
				Body:  template.HTML("<pre><code>" + html.EscapeString(string(data)) + "</code></pre>"),
			}
			app.wirteTemplateToResp(htmltemplate.ContentTemplatePage, http.StatusOK, ctx, w)
		}

		if app.cfg.Debug {
			handler = createDebugHanlder(path.Join(fixedPagesPrefixDir, filePath))
		}

		r.Get("/"+name, handler)
		r.Get("/"+name+"/raw", rawHandler)
		return nil
	})

	if err != nil {
		panic(errors.Wrap(err, "cannot add fixed pages"))
	}
}

// addFileServer sets up a http.FileServer handler to serve static files
func (app *App) addFileServer(r chi.Router, path string) {
	if strings.ContainsAny(path, "{}*") || strings.HasSuffix(path, "/") {
		panic("FileServer does not permit URL parameters or trailing slashes.")
	}

	r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)

	webFS := http.FileServer(app.staticFs)
	fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			app.notFound(w, r)
			return
		}
		webFS.ServeHTTP(w, r)
	})
	r.Method("GET", path+"/{fileName}", http.StripPrefix(path, fileHandler))
	r.Method("GET", "/favicon.ico", http.HandlerFunc(webFS.ServeHTTP))
}
