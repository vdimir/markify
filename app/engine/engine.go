package engine

import (
	"bytes"
	"net/url"
	"time"

	"github.com/vdimir/markify/app/apperr"

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

// SaveDocument process user input data and save document to database
func (eng *DocEngine) SaveDocument(preDoc *UserDocumentData) (DocumentFullSaved, error) {
	doc, err := eng.createDocument(preDoc)
	if err != nil {
		return nil, err
	}
	key, err := eng.docStore.SaveDocument(doc)
	if err != nil {
		return nil, apperr.DBError{Inner: err}
	}
	return &documentWrapper{
		dbDoc: doc,
		key:   key,
	}, nil
}

// LoadDocumentRender return HTML for document with key
func (eng *DocEngine) LoadDocumentRender(key []byte) (DocumentRender, error) {
	var dbDoc = &docstore.MdDocument{}
	err := eng.docStore.LoadDocument(key, docstore.ProjMeta|docstore.ProjRender, dbDoc)
	if err != nil {
		return nil, err
	}
	return &documentWrapper{nil, dbDoc}, err
}

// CreateDocument process user input data and resturs document
func (eng *DocEngine) CreateDocument(preDoc *UserDocumentData) (DocumentFull, error) {
	doc, err := eng.createDocument(preDoc)
	return &documentWrapper{nil, doc}, err
}

func (eng *DocEngine) createDocument(preDoc *UserDocumentData) (*docstore.MdDocument, error) {
	if preDoc.Data == nil || len(preDoc.Data) == 0 {
		err := errors.Errorf("empty data")
		return nil, apperr.WrapfUserError(err, "Data in empty. Type something")
	}

	var textData []byte
	var err error
	var srcURL *url.URL
	if preDoc.IsURL {
		srcURL, err = parseURL(string(preDoc.Data))
		if err != nil {
			return nil, apperr.WrapfUserError(err, "Incorrect URL")
		}
		textData, err = downloadMd(srcURL, eng.fetcher)
		if err != nil {
			return nil, apperr.WrapfUserError(err, "Cannot retrieve data from URL")
		}

	} else {
		textData = preDoc.Data
	}

	curTime := time.Now()
	doc := &docstore.MdDocument{
		MdMeta: docstore.MdMeta{
			CreationTime: curTime.Unix(),
			UpdateTime:   curTime.Unix(),
			MdDocumentParams: docstore.MdDocumentParams{
				EnableShortcodes: preDoc.EnableShortcodes,
			},
		},
		Text: textData,
	}
	if srcURL != nil {
		doc.SrcURL = []byte(srcURL.String())
	}

	err = eng.renderDocument(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (eng *DocEngine) renderDocument(doc *docstore.MdDocument) error {
	ropts := &md.Options{DisableShortcodes: !doc.EnableShortcodes}
	renderHTMLBuf, err := eng.mdrender.Render(doc.Text, ropts)
	if err != nil {
		return errors.Wrap(err, "page render error")
	}
	if renderHTMLBuf.Len() == 0 {
		return apperr.WrapfUserError(errors.New("empty page rendered"), "Empty content!")
	}
	doc.RenderedHTML = renderHTMLBuf.Bytes()

	rawHTMLRender := []byte("<!-- raw HTML omitted -->")
	isEmptyRawHTML := bytes.Compare(bytes.TrimSpace(doc.RenderedHTML), rawHTMLRender) == 0
	if isEmptyRawHTML {
		return apperr.WrapfUserError(errors.New("empty raw HTML page rendered"), "Empty content!")
	}

	return nil
}
