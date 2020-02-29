package app

import (
	"fmt"
	"math/rand"
	"net"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/app/apperr"
	"github.com/vdimir/markify/app/engine"
	"github.com/vdimir/markify/fetch"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../testdata"

func chooseRandomUnusedPort() (port uint16) {
	for i := 0; i < 10; i++ {
		port = 40000 + uint16(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port)); err == nil {
			_ = ln.Close()
			return port
		}
	}
	return 0
}

func createNewTestApp(t *testing.T, tc *TestConfig) (tapp *App, teardown func()) {
	port := chooseRandomUnusedPort()
	require.NotEqual(t, port, 0)

	tmpPath, teardown := testutil.GetTempFolder(t, "test_app")

	tapp, err := NewApp(&Config{
		Debug:        true,
		AssetsPrefix: "../assets",
		DBPath:       tmpPath,
	}, tc)
	assert.NoError(t, err)

	return tapp, teardown
}

func checkURLHash(t *testing.T, urlHash []byte) {
	m, err := regexp.MatchString("[a-zA-Z0-9]+", string(urlHash))
	assert.NoError(t, err)
	assert.Truef(t, m, "unexpected url path %q", string(urlHash))
}

func TestNewURLPage(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	tc := &TestConfig{
		fetcher: fetch.NewMock(),
	}

	tapp, teardown := createNewTestApp(t, tc)
	defer teardown()
	require.NotNil(tapp)

	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))

	tc.fetcher.(*fetch.Mock).SetData("http://foo.bar/page.md", mdData)
	tc.fetcher.(*fetch.Mock).SetData("http://gist.github.com/abc/raw", mdData)
	tc.fetcher.(*fetch.Mock).SetData("file:///home/page.md", mdData)

	userURLInput := func(path string) *engine.UserDocumentData {
		return &engine.UserDocumentData{
			Data:  []byte(path),
			IsURL: true,
		}
	}
	urlHash, err := tapp.saveDocument(userURLInput("http://foo.bar/page.md"))
	assert.NoError(err)
	checkURLHash(t, urlHash)
	doc, err := tapp.getDocument(urlHash)
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile("<h1[a-z\"= ]*>Header</h1>"), doc.HTMLBody())
	assert.Regexp(regexp.MustCompile("<h2[a-z\"= ]*>Subheader</h2>"), doc.HTMLBody())
	assert.Regexp(regexp.MustCompile("Ok"), doc.HTMLBody())

	_, err = tapp.saveDocument(userURLInput("http://goo.gl/page.md"))
	assert.Error(err)
	assert.IsType(apperr.UserError{}, err)

	_, err = tapp.saveDocument(userURLInput("gist.github.com/abc/raw"))
	assert.NoError(err)
	checkURLHash(t, urlHash)

	_, err = tapp.saveDocument(userURLInput("file:///home/page.md"))
	assert.Error(err)
	assert.IsType(apperr.UserError{}, err)

	_, err = tapp.saveDocument(userURLInput("{}dsfsdfa}"))
	assert.Error(err)
	assert.IsType(apperr.UserError{}, err)
}

func TestNewTextPage(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	tc := &TestConfig{
		fetcher: fetch.NewMock(),
	}

	tapp, teardown := createNewTestApp(t, tc)
	defer teardown()
	require.NotNil(tapp)

	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
	key, err := tapp.saveDocument(engine.NewUserDocumentData(mdData))
	assert.NoError(err)

	doc, err := tapp.getDocument(key)
	assert.NoError(err)
	assert.Regexp(regexp.MustCompile("<h1[a-z\"= ]*>Header</h1>"), doc.HTMLBody())
	assert.Regexp(regexp.MustCompile("<h2[a-z\"= ]*>Subheader</h2>"), doc.HTMLBody())
	assert.Regexp(regexp.MustCompile("Ok"), doc.HTMLBody())
}
