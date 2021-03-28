package engine

import (
	"fmt"
	"net/url"
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
