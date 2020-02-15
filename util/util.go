package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
)

// FileWalkerFunc walker function for WalkFiles
type FileWalkerFunc func(data []byte, filePath string) error

// WalkFiles is wrapper around fs.Walk that walks only files exculding direcorines
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
