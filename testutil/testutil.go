package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MustReadData open and read file checking errors
func MustReadData(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return data
}

// GetTempFolder return tmp forder and cleanup func
func GetTempFolder(t require.TestingT, prefix string) (tmpPath string, cleanup func()) {
	tmpPath, err := ioutil.TempDir("", prefix)
	require.NoError(t, err)

	cleanup = func() {
		defer func() {
			assert.NoError(t, os.RemoveAll(tmpPath))
		}()
	}
	return tmpPath, cleanup
}
