package engine

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/vdimir/markify/fetch"
)

func parseURL(rawurl string) (*url.URL, error) {
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

func downloadMd(docURL *url.URL, fetcher fetch.Fetcher) ([]byte, error) {
	dataReader, err := fetcher.Fetch(docURL.String())
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
