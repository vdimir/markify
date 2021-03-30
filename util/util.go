package util

import (
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// FileWalkerFunc walker function for WalkFiles
type FileWalkerFunc func(data []byte, filePath string) error

// WalkFiles is wrapper around fs.Walk that walks only files excluding directories
// and handles file one errors
func WalkFiles(hfs fs.FS, prefixPath string, walker FileWalkerFunc) error {
	walkErr := fs.WalkDir(hfs, prefixPath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrapf(err, "fs walk error")
		}
		if d.IsDir() {
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
