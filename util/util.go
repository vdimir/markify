package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

// FileWalkerFunc walker function for WalkFiles
type FileWalkerFunc func(data []byte, filePath string) error

// WalkFiles is wrapper around fs.Walk that walks only files excluding direcorines
// and handles file one errors
func WalkFiles(hfs http.FileSystem, prefixPath string, walker FileWalkerFunc) error {
	walkErr := fs.Walk(hfs, prefixPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		f, err := hfs.Open(filePath)
		if err != nil {
			return errors.Wrapf(err, "wrong filepath: %q", filePath)
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.Wrapf(err, "read error for %q", filePath)
		}
		subPath := strings.TrimPrefix(filePath, prefixPath+"/")
		err = walker(data, subPath)
		if err != nil {
			return err
		}
		return nil
	})

	return walkErr
}

// GetJSON perform GET request and decode data
func GetJSON(path string, v interface{}) error {
	resp, err := http.Get(path)
	if err != nil {
		return errors.Wrapf(err, "cannot GET %q", path)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("get from %q respond error: %d - %q", path, resp.StatusCode, body)
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return errors.Wrapf(err, "cannot decode data from %q", path)
	}
	return nil
}

// AddRoutePrefix adds prefix to request and handle with h handler
func AddRoutePrefix(prefix string, h http.HandlerFunc) http.Handler {
	if strings.ContainsAny(prefix, "{}") {
		panic("prefix path not permit URL parameters slashes.")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = prefix + r2.URL.Path
		h.ServeHTTP(w, r2)
	})
}
