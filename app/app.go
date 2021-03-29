package app

import (
	"fmt"
	"github.com/vdimir/markify/render"
	"github.com/vdimir/markify/store"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/vdimir/markify/util"
	"github.com/vdimir/markify/view"
)

const defaultURLHashLen = 7

// Config contains application configuration
type Config struct {
	Debug        bool
	AssetsPrefix string
	DBPath       string
	StatusText   string
	UIDSecret    string // secret key to generate user ids
}

type Store interface {
	SetBlob(key string, reader io.Reader, meta map[string]string, ttl time.Duration) error
	GetBlob(key string) (io.Reader, error)
	GetMeta(key string) (map[string]string, error)
	DeleteBlob(key string) error
}

// App provides high level interface to app functions for server
type App struct {
	cfg        *Config
	converter  *render.DocConverter
	blobStore  Store
	uidGen     *util.SignedUIDGenerator
	staticFs   http.FileSystem
	htmlView   view.HTMLPageRender
	httpServer *http.Server
	Addr       string
}

// NewApp create new App instance
func NewApp(cfg *Config) (*App, error) {
	if err := os.MkdirAll(cfg.DBPath, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, err
	}

	var staticFs http.FileSystem
	if cfg.Debug {
		staticFs = newLocalFs(cfg.AssetsPrefix)
	} else {
		staticFs = newStatikFs()
	}

	var htmlView view.HTMLPageRender
	var err error = nil
	if cfg.Debug {
		htmlView, err = view.NewDebugRender(staticFs)
	} else {
		htmlView, err = view.NewRender(staticFs)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error initializing html templates")
	}
	var uidGen *util.SignedUIDGenerator
	if cfg.UIDSecret != "" {
		uidGen = util.NewSignedUIDGenerator([]byte(cfg.UIDSecret))
		cfg.UIDSecret = ""
	}

	blobStore, err := store.NewBoltStorage(path.Join(cfg.DBPath, "data.bdb"))
	if err != nil {
		return nil, errors.Wrap(err, "error initializing storage")
	}
	app := &App{
		cfg:       cfg,
		uidGen:    uidGen,
		converter: render.NewConverter(),
		blobStore: blobStore,
		staticFs:  staticFs,
		htmlView:  htmlView,
	}

	return app, nil
}

var emptyTextRegex = regexp.MustCompile("^\\s*$")

func (app *App) validatePasteRequest(req *CreatePasteRequest) error {
	if !utf8.ValidString(req.Text) {
		return errors.New("got broken input text")

	}
	if emptyTextRegex.MatchString(req.Text) {
		return errors.New("got empty input")
	}
	if app.uidGen == nil || !app.uidGen.Validate([]byte(req.UserToken)) {
		req.UserToken = ""
	}
	if err := app.converter.SupportSyntax(req.Syntax); err != nil {
		return err
	}
	return nil
}

func (app *App) concatTitle(docTitle string, customTitle string) string {
	switch {
	case docTitle != "" && customTitle != "":
		return fmt.Sprintf("%s - %s", customTitle, docTitle)
	case docTitle != "":
		return docTitle
	case customTitle != "":
		return customTitle
	}
	return defaultTitle
}

func (app *App) viewDocument(
	doc *render.Document,
	customTitle string,
	ogURL string,
	w http.ResponseWriter) {
	title := app.concatTitle(doc.Title, customTitle)

	var ogInfo *view.OpenGraphInfo
	if ogURL != "" {
		ogInfo = &view.OpenGraphInfo{
			Title:       title,
			Type:        "article",
			URL:         ogURL,
			Image:       "/public/og-splash.png",
			Description: doc.Preview,
		}
	}

	docView := &view.PageContext{
		Title:  title,
		Body:   template.HTML(doc.Body),
		OgInfo: ogInfo,
	}
	app.viewTemplate(http.StatusOK, docView, w)
}

// Validate request and save data. Returns id of created paste
func (app *App) savePaste(req *CreatePasteRequest) (string, error) {
	docID := util.Base58UID(defaultURLHashLen)

	meta := map[string]string{}
	meta["user"] = req.UserToken
	meta["syntax"] = req.Syntax
	// TODO: meta["create_time"] = ...
	// TODO: meta["ttl"] = ...
	err := app.blobStore.SetBlob(string(docID), strings.NewReader(req.Text), meta, req.Ttl)
	return string(docID), err
}

func (app *App) getDocument(docID string) (*render.Document, error) {
	data, err := app.blobStore.GetBlob(docID)
	if err != nil {
		return nil, err
	}
	meta, err := app.blobStore.GetMeta(docID)
	if err != nil {
		return nil, err
	}
	doc, err := app.converter.Convert(data, meta["syntax"])
	return doc, err
}

func newStatikFs() http.FileSystem {
	log.Printf("[INFO] use assets embedded to binary")
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("[ERROR] no embedded assets loaded, %s", err)
		return nil
	}
	return statikFS
}

func newLocalFs(localPath string) http.FileSystem {
	log.Printf("[INFO] use assets from local directory %q", localPath)
	statikFS := http.Dir(localPath)
	return statikFS
}
