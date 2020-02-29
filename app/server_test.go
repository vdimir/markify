package app_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/app"
	"github.com/vdimir/markify/testutil"
)

func createServer(t *testing.T) (*app.App, func()) {
	port := testutil.ChooseRandomUnusedPort()
	require.NotEqual(t, port, 0)

	tmpPath, tmpFolderClean := testutil.GetTempFolder(t, "test_app")

	tapp, err := app.NewApp(&app.Config{
		Debug:        false,
		AssetsPrefix: "../assets",
		DBPath:       tmpPath,
		StatusText:   `{"status": "ok"}`,
	}, nil)
	require.NoError(t, err)

	go tapp.StartServer("localhost", port)

	err = testutil.WaitForHTTPSServerStart("localhost", port)
	require.NoError(t, err)

	return tapp, func() {
		defer tmpFolderClean()
		defer tapp.Shutdown()
	}
}

func getResp(t *testing.T, path string, expectedStatus int) *http.Response {
	resp, err := http.Get(path)
	assert.NoError(t, err)
	if expectedStatus > 0 {
		assert.Equalf(t, resp.StatusCode, expectedStatus,
			"wrong status code for GET %s", path)
	}
	return resp
}

func TestServerEndpointsExists(t *testing.T) {
	tapp, teardown := createServer(t)
	defer teardown()

	appPath := func(path string) string {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return "http://" + tapp.Addr + path
	}

	paths := []string{
		"/robots.txt",
		"/public/style.css",
		"/about",
		"/info/markdown",
		"/favicon.ico",
	}
	for _, path := range paths {
		_ = getResp(t, appPath(path), http.StatusOK)
	}

	// TODO add more checks
}
