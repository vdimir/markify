package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/vdimir/markify/store/kvstore"
	"go.etcd.io/bbolt"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/vdimir/markify/app/apperr"
	"github.com/vdimir/markify/app/engine"
	"github.com/vdimir/markify/app/view"
	"github.com/vdimir/markify/fetch"
	md "github.com/vdimir/markify/mdrender"
	"github.com/vdimir/markify/util"
)

const defaultURLHashLen = 7

// Config contains application configuration
type Config struct {
	ServerAddrHost string
	ServerPort     uint16
	Debug          bool
	AssetsPrefix   string
	DBPath         string
	StatusText     string
}

// App provides high level interface to app functions for server
type App struct {
	cfg          Config
	engine       *engine.DocEngine
	docIDMapping kvstore.Store
	staticFs     http.FileSystem
	htmlView     view.HTMLPageRender
	httpServer   *http.Server
}

// NewApp create new App instance
func NewApp(cfg Config) (*App, error) {
	docIDMapping, err := kvstore.NewBoltStorage(path.Join(cfg.DBPath, "docid.db"), bbolt.Options{})
	if err != nil {
		return nil, err
	}

	mdRen, err := md.NewRender()
	if err != nil {
		return nil, err
	}
	docEngine := engine.NewDocEngine(cfg.DBPath, mdRen, fetch.NewFetcher())

	var staticFs http.FileSystem
	if cfg.Debug {
		staticFs = newLocalFs(cfg.AssetsPrefix)
	} else {
		staticFs = newStatikFs()
	}

	var htmlView view.HTMLPageRender
	if cfg.Debug {
		htmlView, err = view.NewDebugRender(staticFs)
	} else {
		htmlView, err = view.NewRender(staticFs)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error initializing html templates")
	}

	app := &App{
		cfg:          cfg,
		engine:       docEngine,
		docIDMapping: docIDMapping,
		staticFs:     staticFs,
		htmlView:     htmlView,
	}

	return app, nil
}

func (app *App) concatTitle(docTitle string, customTitle string) string {
	switch {
	case docTitle != "" && customTitle != "":
		return fmt.Sprintf("%s - %s", customTitle, docTitle)
	case docTitle != "":
		return string(docTitle)
	case customTitle != "":
		return customTitle
	}
	return "markify"
}

func (app *App) viewDocument(doc engine.DocumentRender, title string, w http.ResponseWriter) {
	docView := &view.PageContext{
		Title: app.concatTitle(doc.Title(), title),
		Body:  doc.HTMLBody(),
	}
	app.viewTemplate(http.StatusOK, docView, w)
}

func (app *App) viewRawDocument(doc engine.DocumentText, title string, w http.ResponseWriter) {
	docView := &view.PageContext{
		Title: app.concatTitle(doc.Title(), title),
		Body:  template.HTML(fmt.Sprintf("<pre><code>%s</code></pre>", doc.MdText())),
	}
	app.viewTemplate(http.StatusOK, docView, w)
}

func (app *App) createDocID(doc engine.DocumentSaved) ([]byte, error) {
	newID := util.Base58UID(defaultURLHashLen)
	err := app.docIDMapping.Save(newID, doc.Key())
	if err != nil {
		return nil, apperr.DBError{err}
	}
	return newID, nil
}

func (app *App) saveDocument(preDoc *engine.UserDocumentData) ([]byte, error) {
	doc, err := app.engine.SaveDocument(preDoc)
	if err != nil {
		return nil, err
	}
	docID, err := app.createDocID(doc)
	if err != nil {
		return nil, err
	}
	return docID, nil
}

func (app *App) getDocument(docID []byte) (engine.DocumentRender, error) {
	key, err := app.docIDMapping.Load(docID)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, nil
	}
	return app.engine.LoadDocumentRender(key)
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
