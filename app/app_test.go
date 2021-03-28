package app

import (
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/app/engine"
	"github.com/vdimir/markify/fetch"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../testdata"

func createNewTestApp(t *testing.T, tc *TestConfig) (tapp *App, teardown func()) {
	tmpPath, tmpFolderClean := testutil.GetTempFolder(t, "test_app")

	tapp, err := NewApp(&Config{
		Debug:        false,
		AssetsPrefix: "../assets",
		DBPath:       tmpPath,
	}, tc)
	assert.NoError(t, err)

	return tapp, tmpFolderClean
}

func checkURLHash(t *testing.T, urlHash []byte) {
	m, err := regexp.MatchString("[a-zA-Z0-9]+", string(urlHash))
	assert.NoError(t, err)
	assert.Truef(t, m, "unexpected url path %q", string(urlHash))
}

func TestNewMarkdownPage(t *testing.T) {
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

	{
		key, err := tapp.saveDocument(engine.NewUserDocumentData([]byte("   ")))
		assert.Error(err)
		assert.Nil(key)
	}
	{
		key, err := tapp.saveDocument(engine.NewUserDocumentData([]byte("<p>abc</p>\n<div>def</div>")))
		assert.Error(err)
		assert.Nil(key)
	}
}
