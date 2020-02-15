package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/vdimir/markify/fetch"
	"github.com/vdimir/markify/render/htmltemplate"
	"github.com/vdimir/markify/render/md"
	"github.com/vdimir/markify/store"
	"github.com/vdimir/markify/util"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	bolt "go.etcd.io/bbolt"
)

const defaultURLHashLen = 7

// Config contains application configuration
type Config struct {
	ServerAddrHost string
	ServerPort     uint16
	Debug          bool
	AssetsPrefix   string
	PageCachePath  string
	KeyStorePath   string
	MdTextPath     string
}

// App contains application parts
type App struct {
	cfg           Config
	pageCache     store.Store
	keyStore      store.Store
	rawTextStore  store.KeyStore
	render        *md.Render
	fetcher       fetch.Fetcher
	staticFs      http.FileSystem
	htmlTplRender htmltemplate.HTMLPageRender
	httpServer    *http.Server
}

// NewApp create new App instance
func NewApp(cfg Config) (*App, error) {
	keyStore, err := store.NewBoltStorage(cfg.KeyStorePath, bolt.Options{})
	if err != nil {
		return nil, err
	}

	pageStore, err := store.NewBoltStorage(cfg.PageCachePath, bolt.Options{})
	if err != nil {
		return nil, err
	}

	rawTextStore, err := store.NewBoltStorage(cfg.MdTextPath, bolt.Options{})
	if err != nil {
		return nil, err
	}

	var localPath *string
	if cfg.Debug {
		localPath = new(string)
		*localPath = cfg.AssetsPrefix
	}

	htmlRen, err := md.NewRender()
	if err != nil {
		return nil, err
	}

	staticFs := newStaticFs(localPath)

	var htmlTplRender htmltemplate.HTMLPageRender
	if cfg.Debug {
		htmlTplRender, err = htmltemplate.NewDebugRender(staticFs)
	} else {
		htmlTplRender, err = htmltemplate.NewRender(staticFs)
	}
	if err != nil {
		return nil, errors.Wrap(err, "error initalizing html templates")
	}

	app := &App{
		cfg:           cfg,
		pageCache:     pageStore,
		keyStore:      keyStore,
		render:        htmlRen,
		fetcher:       fetch.NewFetcher(),
		staticFs:      staticFs,
		htmlTplRender: htmlTplRender,
		rawTextStore:  rawTextStore,
	}
	return app, nil
}

func (app *App) saveRenderToCache(rawMdData []byte, urlHash []byte, renderOpt *md.Options) error {

	htmlBuf, err := app.render.Render(rawMdData, renderOpt)
	if err != nil {
		return errors.Wrap(err, "page render error")
	}

	err = app.pageCache.Save(urlHash, htmlBuf.Bytes())
	if err != nil {
		return DBError{err}
	}
	return err
}

// getRawFromURL download raw markdown from url
func (app *App) getRawFromURL(pageURL *url.URL) ([]byte, error) {
	dataReader, err := app.fetcher.Fetch(pageURL.String())
	if err != nil {
		return nil, err
	}
	defer dataReader.Close()

	rawMdData, err := ioutil.ReadAll(dataReader)
	if err != nil {
		return nil, err
	}
	return rawMdData, nil
}

func (app *App) newShortEncode(data []byte) ([]byte, error) {
	urlHash, _ := util.BaseHashEncode(data, defaultURLHashLen)
	if t, err := app.keyStore.Timestamp(urlHash); t != 0 || err != nil {
		if err != nil {
			return nil, DBError{err}
		}
		return nil, fmt.Errorf("key collision: %q exists  in db", urlHash)
	}
	return urlHash, nil
}

// addPageByURL download raw markdown from url and save rendered page
func (app *App) addPageByURL(params formParams) ([]byte, error) {
	pageURL, err := app.parseURL(string(params.TextData))
	if err != nil {
		return nil, UserError{err}
	}

	rawMdData, err := app.getRawFromURL(pageURL)
	if err != nil {
		return nil, UserError{err}
	}

	urlHash, err := app.newShortEncode([]byte(pageURL.String()))
	if err != nil {
		return nil, err
	}

	opt := &md.Options{BaseURL: pageURL, DisableShortcodes: !params.EnableShortcodes}
	err = app.saveRenderToCache(rawMdData, urlHash, opt)
	if err != nil {
		return nil, err
	}

	err = app.keyStore.Save(urlHash, []byte(pageURL.String()))
	if err != nil {
		return nil, DBError{err}
	}

	return urlHash, nil
}

func (app *App) addPageByText(formData formParams) ([]byte, error) {
	key, err := app.rawTextStore.NewKey(formData.TextData)
	if err != nil {
		return nil, DBError{err}
	}

	urlHash, err := app.newShortEncode(key)
	if err != nil {
		return nil, err
	}

	ropt := &md.Options{DisableShortcodes: !formData.EnableShortcodes}
	err = app.saveRenderToCache(formData.TextData, urlHash, ropt)
	if err != nil {
		return nil, err
	}

	localURL := bytes.NewBuffer([]byte("local://"))
	localURL.Write(key)
	err = app.keyStore.Save(urlHash, localURL.Bytes())
	if err != nil {
		return nil, DBError{err}
	}

	return urlHash, nil
}

func (app *App) parseURL(rawurl string) (*url.URL, error) {
	pageURL, err := url.Parse(rawurl)
	if err == nil && pageURL.Scheme == "" {
		pageURL, err = url.ParseRequestURI("http://" + rawurl)
	}
	if err != nil {
		return nil, err
	}

	if pageURL.Scheme != "http" && pageURL.Scheme != "https" {
		return nil, fmt.Errorf("Ursupported scheme for %v", pageURL)
	}
	return pageURL, err
}

func newStaticFs(localPath *string) http.FileSystem {
	var statikFS http.FileSystem

	if localPath != nil {
		log.Printf("[INFO] use assets from local directory %q", *localPath)
		statikFS = http.Dir(*localPath)
	} else {
		log.Printf("[INFO] use assets embeded to binary")
		var err error
		statikFS, err = fs.New()
		if err != nil {
			log.Fatalf("[ERROR] no embedded assets loaded, %s", err)
			return nil
		}
	}
	return statikFS
}
