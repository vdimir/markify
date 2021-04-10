package util

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalkFilesHttpDir(t *testing.T) {
	fs := os.DirFS("../view")
	filesNames := map[string]bool{}
	err := WalkFiles(fs, "template", func(data []byte, filePath string) error {
		filesNames[filePath] = true
		assert.True(t, data != nil && len(data) > 0)
		return nil
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(filesNames), 5)
	// Check some files exists
	assert.Contains(t, filesNames, "page.html")
	assert.Contains(t, filesNames, "partial/footer.html")
}

func TestAddRoutePrefix(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Path", r.URL.Path)
	})

	testCases := []string{
		"/bar", "/bar/biz",
		"/bar/", "/bar/biz/",
	}
	ts := httptest.NewServer(AddRoutePrefix("/foo/", h))
	c := ts.Client()
	for _, tc := range testCases {
		res, err := c.Get(ts.URL + tc)
		assert.NoError(t, err)
		assert.Equal(t, res.Header.Get("X-Path"), "/foo/"+tc)
	}
}
