package fetch

import (
	"fmt"
	"io"
	"mime"
	"net/http"
)

// Fetcher download data from source
type Fetcher interface {
	Fetch(url string) (io.ReadCloser, error)
}

// SimpleFetcher download data from source and checks content type
type SimpleFetcher struct {
	contetTypes map[string]struct{}
}

// NewFetcher create new Fetcher
func NewFetcher() Fetcher {
	contetTypes := map[string]struct{}{
		"text/plain": struct{}{},
	}
	return SimpleFetcher{
		contetTypes: contetTypes,
	}
}

// Fetch retrieve data from url
func (f SimpleFetcher) Fetch(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if _, ok := f.contetTypes[contentType]; !ok {
		return nil, fmt.Errorf("unsupported content type %s", contentType)
	}
	return resp.Body, err
}
