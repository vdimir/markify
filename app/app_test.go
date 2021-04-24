package app

import (
	"fmt"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../testdata"

func createNewTestApp(t *testing.T) (tapp *App, teardown func()) {
	tmpPath, tmpFolderClean := testutil.GetTempFolder(t, "test_app")

	tapp, err := NewApp(&Config{
		Debug:        false,
		AssetsPrefix: "assets",
		StorageSpec:  fmt.Sprintf("local:%s", tmpPath),
		StatusText:   `{"status": "ok"}`,
	})
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

	tapp, teardown := createNewTestApp(t)
	defer teardown()
	require.NotNil(tapp)

	mdData := testutil.MustReadData(t, path.Join(testDataPath, "page.md"))
	key, err := tapp.savePaste(&CreatePasteRequest{Text: string(mdData), Syntax: "markdown"})
	assert.NoError(err)
	{
		doc, err := tapp.getDocument(key)
		require.NoError(err)
		require.NotNil(doc)
		assert.Regexp(regexp.MustCompile("<h1[a-z\"= ]*>Header</h1>"), doc.Body)
		assert.Regexp(regexp.MustCompile("<h2[a-z\"= ]*>Subheader</h2>"), doc.Body)
		assert.Regexp(regexp.MustCompile("Ok"), doc.Body)
	}
	{
		unexistingKey := "__deadbeef__"
		doc, err := tapp.getDocument(unexistingKey)
		require.NoError(err)
		require.Nil(doc)
	}
}
