package engine

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/vdimir/markify/fetch"
	md "github.com/vdimir/markify/mdrender"
	"github.com/vdimir/markify/store/docstore"
)

// DocEngine deals with documents: rendering, saving, etc
type DocEngine struct {
	docStore docstore.DocStore
	mdrender *md.Render
	fetcher  fetch.Fetcher
}

// NewDocEngine creates new DocEngine
func NewDocEngine(dbPath string, mdrender *md.Render, fetcher fetch.Fetcher) *DocEngine {
	return &DocEngine{
		docStore: docstore.NewBoltDocStore(dbPath),
		mdrender: mdrender,
		fetcher:  fetcher,
	}
}

func (eng *DocEngine) SaveDocument(preDoc *UserDocumentData) (DocumentFullSaved, error) {
	doc, err := eng.createDocument(preDoc)
	if err != nil {
		return nil, err
	}
	key, err := eng.docStore.SaveDocument(doc)
	if err != nil {
		return nil, err
	}
	return &documentWrapper{
		dbDoc: doc,
		key:   key,
	}, nil
}

func (eng *DocEngine) LoadDocumentRender(key []byte) (DocumentRender, error) {
	var dbDoc *docstore.MdDocument
	err := eng.docStore.LoadDocument(key, docstore.ProjMeta|docstore.ProjRender, dbDoc)
	if err != nil {
		return nil, err
	}
	return &documentWrapper{nil, dbDoc}, err
}

func (eng *DocEngine) CreateDocument(preDoc *UserDocumentData) (DocumentFull, error) {
	doc, err := eng.createDocument(preDoc)
	return &documentWrapper{nil, doc}, err
}

func (eng *DocEngine) createDocument(preDoc *UserDocumentData) (*docstore.MdDocument, error) {
	var textData []byte
	var err error
	var srcURL *url.URL
	if preDoc.IsURL {
		srcURL, textData, err = eng.downloadDocument(string(preDoc.Data))
	} else {
		textData = preDoc.Data
	}

	curTime := time.Now()
	doc := &docstore.MdDocument{
		SrcURL:       srcURL,
		Text:         textData,
		CreationTime: curTime,
		UpdateTime:   curTime,
		Params: docstore.MdDocumentParams{
			EnableShortcodes: preDoc.EnableShortcodes,
		},
	}

	err = eng.renderDocument(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (eng *DocEngine) downloadDocument(rawurl string) (*url.URL, []byte, error) {
	srcURL, err := parseURL(rawurl)
	if err != nil {
		return nil, nil, err
	}

	var data []byte
	data, err = downloadMd(srcURL, eng.fetcher)
	if err != nil {
		return nil, nil, err
	}
	return srcURL, data, nil
}

func (eng *DocEngine) renderDocument(doc *docstore.MdDocument) error {
	ropts := &md.Options{DisableShortcodes: !doc.Params.EnableShortcodes}
	renderHTMLBuf, err := eng.mdrender.Render(doc.Text, ropts)
	if err != nil {
		return errors.Wrap(err, "page render error")
	}
	doc.RenderedHTML = renderHTMLBuf.Bytes()
	return nil
}
