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

	type TestCase struct {
		text          string
		expectedTitle string
		expectedDesc  string
		n             int
	}

	testCases := []TestCase{
		{
			text:          "# Page title",
			expectedTitle: "^Page title$",
			expectedDesc:  "^$",
		},
		{
			text:          "Page title",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Page title$",
		},
		{
			text: "# Page title" + "\n" +
				"## L2 Header" + "\n" +
				"some content" + "\n" +
				"## Other L2 Header" + "\n" +
				"some other content" + "\n" +
				"# H1 Header" + "\n" +
				"xxx" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^some content$",
		},
		{
			text: "## Page title" + "\n" +
				"### L3 Header" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^$",
		},
		{
			text: "Some content before header" + "\n" +
				"## Page title" + "\n" +
				"### L3 Header" + "\n" +
				"Other content" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Some content before header$",
		},
		{
			text: "Some content before header" + "\n" +
				"### Page title" + "\n" +
				"### L3 Header" + "\n" +
				"Other content" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Some content before header$",
		},
		{
			text: "Some content before header" + "\n" +
				"#### Page title" + "\n" +
				"# L1 Header" + "\n" +
				"Other content" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Some content before header$",
		},
		{
			text: "\n\n" +
				"Page title" + "\n\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Page title$",
		},
		{
			text: "\n\n" +
				"Page\ntitle" + "\n\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n",
			expectedTitle: "^Page title$",
			expectedDesc:  "^Page title$",
		},
		{
			text: "\n\n" +
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit, " +
				"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. " +
				"Ut enim ad minim veniam, quis nostrud exercitation ullamco" +
				"laboris nisi ut aliquip ex ea commodo consequat." + "\n" +
				"No headers at all" + "\n\n" +
				"Other content" + "\n\n",
			expectedTitle: "^Lorem ipsum dolor sit amet.{0,40}…$",
			expectedDesc:  "^Lorem ipsum dolor sit amet[A-Za-z \\.,]{0,250}$",
		},
	}

	checkDoc := func(doc engine.Document, testCase TestCase) {
		titleRe := regexp.MustCompile(testCase.expectedTitle)
		assert.Regexpf(t, titleRe, doc.Title(), "Error in test case #%d", testCase.n)
		descRe := regexp.MustCompile(testCase.expectedDesc)
		assert.Regexpf(t, descRe, doc.Description(), "Error in test case #%d", testCase.n)
	}

	for n, testCase := range testCases {
		testCase.n = n + 1
		doc, err := docEng.CreateDocument(engine.NewUserDocumentData([]byte(testCase.text)))
		assert.NoErrorf(t, err, "Error in test case #%d", testCase.n)
		checkDoc(doc, testCase)
	}

	for n, testCase := range testCases {
		testCase.n = n + 1
		doc, err := docEng.SaveDocument(engine.NewUserDocumentData([]byte(testCase.text)))
		assert.NoErrorf(t, err, "Error in test case #%d", testCase.n)
		checkDoc(doc, testCase)
	}
}
