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

func createNewTestApp(t *testing.T) (tapp *App, teardown func()) {
	port := chooseRandomUnusedPort()
	require.NotEqual(t, port, 0)

	tmpPath, teardown := testutil.GetTempFolder(t, "test_app")

	tapp, err := NewApp(Config{
		ServerAddrHost: "localhost",
		ServerPort:     port,
		Debug:          true,
		AssetsPrefix:   "../assets",
		PageCachePath:  path.Join(tmpPath, "cache.db"),
		URLHashPath:    path.Join(tmpPath, "keys.db"),
		MdTextPath:     path.Join(tmpPath, "mdtext.db"),
	})
	assert.NoError(t, err)
	tapp.fetcher = fetch.NewMock()

	return tapp, teardown
}

func checkURLHash(t *testing.T, urlHash []byte) {
	m, err := regexp.MatchString("[a-zA-Z0-9]+", string(urlHash))
	assert.NoError(t, err)
	assert.Truef(t, m, "unexpected url hash %q", string(urlHash))
}

func TestNewPage(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	tapp, teardown := createNewTestApp(t)
	defer teardown()
	require.NotNil(tapp)

	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))

	tapp.fetcher.(*fetch.Mock).SetData("http://foo.bar/page.md", mdData)
	tapp.fetcher.(*fetch.Mock).SetData("http://gist.github.com/abc/raw", mdData)
	tapp.fetcher.(*fetch.Mock).SetData("file:///home/page.md", mdData)

	params := func(path string) formParams {
		return formParams{TextData: []byte(path)}
	}
	urlHash, err := tapp.addPageByURL(params("http://foo.bar/page.md"))
	assert.NoError(err)
	checkURLHash(t, urlHash)

	_, err = tapp.addPageByURL(params("http://goo.gl/page.md"))
	assert.Error(err)
	assert.IsType(UserError{}, err)

	_, err = tapp.addPageByURL(params("gist.github.com/abc/raw"))
	assert.NoError(err)
	checkURLHash(t, urlHash)

	_, err = tapp.addPageByURL(params("file:///home/page.md"))
	assert.Error(err)
	assert.IsType(UserError{}, err)

	_, err = tapp.addPageByURL(params("{}dsfsdfa}"))
	assert.Error(err)
	assert.IsType(UserError{}, err)

	data, err := tapp.urlHashStore.Load(urlHash)
	assert.NoError(err)
	assert.Equal(data, []byte("http://foo.bar/page.md"))

	data, err = tapp.pageCache.Load(urlHash)
	assert.NoError(err)
	assert.Greater(len(data), 10)
}
