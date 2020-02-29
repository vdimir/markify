package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const appHostURL = "https://test.markify.dev"

func appPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return appHostURL + path
	}
	return appHostURL + "/" + path
}

func checkHandlerResp(t *testing.T, path string, handler http.HandlerFunc, checkFunc func(*http.Response)) {
	var req *http.Request
	var w *httptest.ResponseRecorder

	req = httptest.NewRequest("GET", appPath(path), nil)
	w = httptest.NewRecorder()
	handler(w, req)
	checkFunc(w.Result())
}

func TestHandlersDirectCall(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	tapp, teardown := createNewTestApp(t, nil)
	defer teardown()
	require.NotNil(tapp)

	checkStatusOk := func(r *http.Response) {
		assert.Equal(http.StatusOK, r.StatusCode)
	}

	checkHandlerResp(t, "", tapp.handlePageIndex, checkStatusOk)
	checkHandlerResp(t, "ping", tapp.handlePing, func(r *http.Response) {
		assert.Equal(http.StatusOK, r.StatusCode)
		assert.Contains(r.Header.Get("Content-Type"), "application/json")
	})

	checkHandlerResp(t, "robots.txt", tapp.handleRobotsTxt, checkStatusOk)
	checkHandlerResp(t, "url", tapp.handlePageInputURL, checkStatusOk)
	checkHandlerResp(t, "compose", tapp.handlePageTextInput, checkStatusOk)
}

func TestRoutes(t *testing.T) {
	tapp, teardown := createNewTestApp(t, nil)
	defer teardown()

	ts := httptest.NewServer(tapp.Routes())
	defer ts.Close()

	getResp := func(path string) *http.Response {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		resp, err := ts.Client().Get(ts.URL + path)
		assert.NoError(t, err)
		return resp
	}

	var resp *http.Response

	resp = getResp("/")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// TODO add more checks
}
