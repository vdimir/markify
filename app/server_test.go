package app_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/app"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../testdata"

func createServer(t *testing.T, customCfg func(*app.Config)) (*app.App, func()) {
	port := testutil.ChooseRandomUnusedPort()
	require.NotEqual(t, port, 0)

	tmpPath, tmpFolderClean := testutil.GetTempFolder(t, "test_app")

	cfg := &app.Config{
		Debug:        false,
		AssetsPrefix: "../assets",
		DBPath:       tmpPath,
		StatusText:   `{"status": "ok"}`,
	}
	if customCfg != nil {
		customCfg(cfg)
	}
	tapp, err := app.NewApp(cfg, nil)
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

func createPathHelper(host string) func(path string) string {
	return func(path string) string {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		return "http://" + host + path
	}
}

func mustReadAll(r io.Reader) string {
	data, _ := ioutil.ReadAll(r)
	return string(data)
}

func TestServerEndpointsExists(t *testing.T) {
	tapp, teardown := createServer(t, nil)
	defer teardown()

	appPath := createPathHelper(tapp.Addr)

	paths := []string{
		"/robots.txt",
		"/public/style.css",
		"/about",
		"/info/markdown",
		"/favicon.ico",
		"/compose",
		"/link",
	}
	for _, path := range paths {
		_ = getResp(t, appPath(path), http.StatusOK)
	}

	pathsNotFound := []string{
		"/public/", "/public",
		"/assets/", "/assets",
		"/assets/template/page.html",
		"/info/",
	}
	for _, path := range pathsNotFound {
		_ = getResp(t, appPath(path), http.StatusNotFound)
	}
}

func TestServerDebugMode(t *testing.T) {
	tappDebug, teardown := createServer(t, func(c *app.Config) {
		c.Debug = true
	})
	defer teardown()

	tappNoDebug, teardown := createServer(t, func(c *app.Config) {
		c.Debug = false
	})
	defer teardown()

	appDebugPath := createPathHelper(tappDebug.Addr)
	appNoDebugPath := createPathHelper(tappNoDebug.Addr)

	paths := []string{
		"/robots.txt",
		"/public/style.css",
		"/about",
		"/info/markdown",
		"/",
	}
	for _, path := range paths {
		respDebug := getResp(t, appDebugPath(path), http.StatusOK)
		respNoDebug := getResp(t, appNoDebugPath(path), http.StatusOK)

		assert.Equal(t, respDebug.StatusCode, respNoDebug.StatusCode)

		dataDebug, err := ioutil.ReadAll(respDebug.Body)
		assert.NoError(t, err)

		dataNoDebug, err := ioutil.ReadAll(respNoDebug.Body)
		assert.NoError(t, err)

		assert.Equal(t, dataDebug, dataNoDebug)
	}
}

func TestServerCreatePage(t *testing.T) {
	tapp, teardown := createServer(t, nil)
	defer teardown()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
		w.Write(mdData)
	}))

	appPath := createPathHelper(tapp.Addr)

	composePage := func(data string) *http.Response {
		formData := url.Values{
			"data": {data},
		}
		resp, err := http.PostForm(appPath("/compose"), formData)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		return resp
	}

	linkPage := func(path string) *http.Response {
		formData := url.Values{
			"data": {path},
			"type": {"url"},
		}
		resp, err := http.PostForm(appPath("/link"), formData)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		return resp
	}

	var resp *http.Response
	{
		resp, _ = http.PostForm(appPath("/compose"), url.Values{"data": {""}})
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
	savedPagesCases := map[string]*regexp.Regexp{}
	{
		resp = composePage("foo")
		assert.Contains(t, mustReadAll(resp.Body), "<p>foo</p>")
		assert.True(t, strings.HasPrefix(resp.Request.URL.Path, "/p/"))
		savedPagesCases[resp.Request.URL.Path] = regexp.MustCompile("foo")
	}
	{
		resp = linkPage(ts.URL)
		respData := mustReadAll(resp.Body)
		assert.Regexp(t, regexp.MustCompile("<h1[a-z\"= ]*>Header</h1>"), respData)
		assert.True(t, strings.HasPrefix(resp.Request.URL.Path, "/p/"))
		savedPagesCases[resp.Request.URL.Path] = regexp.MustCompile("Ok")
	}
	{
		mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
		resp = composePage(string(mdData))
		respData := mustReadAll(resp.Body)
		assert.Regexp(t, regexp.MustCompile("<h1[a-z\"= ]*>Header</h1>"), respData)
		assert.Regexp(t, regexp.MustCompile("<h2[a-z\"= ]*>Subheader</h2>"), respData)
		assert.Regexp(t, regexp.MustCompile("Ok"), respData)
		assert.True(t, strings.HasPrefix(resp.Request.URL.Path, "/p/"))
		savedPagesCases[resp.Request.URL.Path] = regexp.MustCompile("Ok")
	}

	for path, expected := range savedPagesCases {
		resp, err := http.Get(appPath(path))
		assert.NoError(t, err)
		respData := mustReadAll(resp.Body)
		assert.Regexp(t, expected, respData)
	}
}
