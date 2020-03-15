package engine_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vdimir/markify/app/engine"
	"github.com/vdimir/markify/fetch"
	md "github.com/vdimir/markify/mdrender"
	"github.com/vdimir/markify/testutil"
)

const testDataPath = "../../testdata"

func initDocEngine(t *testing.T) (*engine.DocEngine, func()) {
	tmpPath, tmpFolderClean := testutil.GetTempFolder(t, "test_app")

	mdren, err := md.NewRender()
	assert.NoError(t, err)

	teng := engine.NewDocEngine(tmpPath, mdren, &fetch.Mock{})
	assert.NotNil(t, teng)

	return teng, tmpFolderClean
}

func TestDocTitle(t *testing.T) {

	docEng, teardown := initDocEngine(t)
	defer teardown()

	testCases := []struct {
		text     string
		expected string
	}{
		{
			text:     "# Page title",
			expected: "^Page title$",
		},
		{
			text:     "Page title",
			expected: "^Page title$",
		},
		// {
		// 	text:     "\n\n",
		// 	expected: "^$",
		// },
		{
			text: "# Page title" + "\n" +
				"## L2 Header" + "\n" +
				"some content" + "\n" +
				"## Other L2 Header" + "\n" +
				"some other content" + "\n" +
				"# H1 Header" + "\n" +
				"xxx" + "\n",
			expected: "^Page title$",
		},
		{
			text: "## Page title" + "\n" +
				"### L3 Header" + "\n",
			expected: "^Page title$",
		},
		{
			text: "Some content before header" + "\n" +
				"## Page title" + "\n" +
				"### L3 Header" + "\n" +
				"Other content" + "\n",
			expected: "^Page title$",
		},
		{
			text: "Some content before header" + "\n" +
				"### Page title" + "\n" +
				"### L3 Header" + "\n" +
				"Other content" + "\n",
			expected: "^Page title$",
		},
		{
			text: "Some content before header" + "\n" +
				"#### Page title" + "\n" +
				"# L1 Header" + "\n" +
				"Other content" + "\n",
			expected: "^Page title$",
		},
		{
			text: "\n\n" +
				"Page title" + "\n\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n",
			expected: "^Page title$",
		},
		{
			text: "\n\n" +
				"Page\ntitle" + "\n\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n",
			expected: "^Page title$",
		},
		{
			text: "\n\n" +
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit, " +
				"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
				"Ut enim ad minim veniam, quis nostrud exercitation ullamco" +
				"laboris nisi ut aliquip ex ea commodo consequat." + "\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n\n",
			expected: "^Lorem ipsum dolor sit amet.{0,40}…$",
		},
	}

	for n, testCase := range testCases {
		doc, err := docEng.CreateDocument(engine.NewUserDocumentData([]byte(testCase.text)))
		testCaseN := n + 1
		assert.NoErrorf(t, err, "Error in CreateDocument #%d", testCaseN)
		expRe := regexp.MustCompile(testCase.expected)
		assert.Regexpf(t, expRe, doc.Title(), "Error in CreateDocument #%d", testCaseN)
	}

	for n, testCase := range testCases {
		doc, err := docEng.SaveDocument(engine.NewUserDocumentData([]byte(testCase.text)))
		testCaseN := n + 1
		assert.NoErrorf(t, err, "Error in SaveDocument #%d", testCaseN)
		expRe := regexp.MustCompile(testCase.expected)
		assert.Regexpf(t, expRe, doc.Title(), "Error in SaveDocument #%d", testCaseN)
	}
}
