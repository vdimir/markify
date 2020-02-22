package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalkFilesHttpDir(t *testing.T) {
	statikFs := http.Dir("../assets")
	filesNames := map[string]bool{}
	err := WalkFiles(statikFs, "/template", func(data []byte, filePath string) error {
		filesNames[filePath] = true
		assert.True(t, data != nil && len(data) > 0)
		return nil
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(filesNames), 5)
	// Check some files exists
	assert.Contains(t, filesNames, "page.html")
	assert.Contains(t, filesNames, "partial/footer.html")
}
