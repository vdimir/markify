package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/vdimir/markify/render"
	"github.com/vdimir/markify/store"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/vdimir/markify/util"
	"github.com/vdimir/markify/view"
)

const defaultURLHashLen = 7

// Config contains application configuration
type Config struct {
	Debug        bool
	AssetsPrefix string
	StorageSpec       string
	StatusText   string
	UIDSecret    string // secret key to generate user ids
}

type Store interface {
	SetBlob(key string, reader io.Reader, meta map[string]string, ttl time.Duration) error
	GetBlob(key string) (io.Reader, map[string]string, error)
	DeleteBlob(key string) error
}

// App provides high level interface to app functions for server
type App struct {
	cfg        *Config
	converter  *render.DocConverter
	blobStore  Store
	uidGen     *util.SignedUIDGenerator
	staticFs   fs.FS
	htmlView   view.HTMLPageView
	httpServer *http.Server
	Addr       string
}

type Document struct {
	render.Document
	CreateTime time.Time
}

// NewApp create new App instance
func NewApp(cfg *Config) (*App, error) {
	var staticFs fs.FS
	if cfg.Debug {
		staticFs = newLocalFs(cfg.AssetsPrefix)
	} else {
		staticFs = newEmbeddedFs()
	}

	templatePath := ""
	if cfg.Debug {
		templatePath = "view/template"
	}
	htmlView, err := view.NewView(templatePath)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing html templates")
	}
	var uidGen *util.SignedUIDGenerator
	if cfg.UIDSecret != "" {
		uidGen = util.NewSignedUIDGenerator([]byte(cfg.UIDSecret))
		cfg.UIDSecret = ""
	}

	blobStore, err := createStorage(cfg.StorageSpec)
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
	doc *Document,
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
	if !doc.CreateTime.IsZero() {
		docView.CreateTime = doc.CreateTime.Format("Jan 2 15:04:05 2006 MST")
	}
	app.viewTemplate(http.StatusOK, docView, w)
}

// Validate request and save data. Returns id of created paste
func (app *App) savePaste(req *CreatePasteRequest) (string, error) {
	docID := util.Base58UID(defaultURLHashLen)

	meta := map[string]string{}
	meta["user"] = req.UserToken
	meta["syntax"] = req.Syntax
	timeStr, err := time.Now().UTC().MarshalText()
	if err != nil {
		return "", err
	}
	meta["create_time"] = string(timeStr)
	meta["ttl"] = req.Ttl.String()
	err = app.blobStore.SetBlob(string(docID), strings.NewReader(req.Text), meta, req.Ttl)
	return string(docID), err
}

func (app *App) getDocument(docID string) (*Document, error) {
	data, meta, err := app.blobStore.GetBlob(docID)
	if err != nil {
		return nil, errors.Wrapf(err, "can't get data")
	}

	rdoc, err := app.converter.Convert(data, meta["syntax"])
	if err != nil {
		return nil, err
	}
	createTime := time.Time{}
	err = createTime.UnmarshalText([]byte(meta["create_time"]))
	if err != nil {
		log.Printf("[ERROR] can't parse time from metadata: %q: %s", meta["create_time"], err.Error())
	}
	doc := &Document{Document: *rdoc, CreateTime: createTime}
	return doc, nil
}

func createStorage(storageSpec string) (Store, error){
	typeAndOptions := strings.SplitN(storageSpec, ":", 2)
	if len(typeAndOptions) != 2 {
		return nil, errors.Errorf("error parse storage specification %q", storageSpec)
	}
	if typeAndOptions[0] == "local" {
		dbFile := path.Join(typeAndOptions[1], "data.bdb")
		log.Printf("[INFO] creating local storage, data file %q", dbFile)
		return store.NewBoltStorage(dbFile)
	}
	if typeAndOptions[0] == "s3" {
		s3conf := store.S3Config{}
		err := json.Unmarshal([]byte(typeAndOptions[1]), &s3conf)
		if err != nil {
			return nil, errors.Wrapf(err, "error parse s3 config")
		}
		log.Printf("[INFO] using s3 storage, endpoint %q, bucket %q", s3conf.Endpoint, s3conf.Bucket)
		return store.NewS3Storage(s3conf)
	}
	return nil, errors.Errorf("unknown storage type %q", typeAndOptions[0])
}

//go:embed assets/*
var embeddedStaticFS embed.FS

func newEmbeddedFs() fs.FS {
	log.Printf("[INFO] use assets embedded to binary")

	embFs, err := fs.Sub(embeddedStaticFS, "assets")
	if err != nil {
		log.Fatal("[ERROR] can't load static files from assets")
	}
	return embFs
}

func newLocalFs(localPath string) fs.FS {
	log.Printf("[INFO] use assets from local directory %q", localPath)
	localFS := os.DirFS(localPath)
	return localFS
}
