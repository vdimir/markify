package docstore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vdimir/markify/testutil"
)

func now() time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	t, err := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (PST)")
	if err != nil {
		panic(err)
	}
	return t
}

func TestBoltDocstoreSaveEmpty(t *testing.T) {
	tmpPath, cleanup := testutil.GetTempFolder(t, "bolt_store")
	defer cleanup()

	s := NewBoltDocStore(tmpPath)
	var err error
	_, err = s.SaveDocument(&MdDocument{})
	assert.Error(t, err)

	_, err = s.SaveDocument(&MdDocument{
		MdMeta: MdMeta{
			CreationTime: now().Unix(),
			UpdateTime:   now().Unix(),
		},
	})
	assert.Error(t, err)
}

func TestBoltDocstoreSaveLoad(t *testing.T) {
	tmpPath, cleanup := testutil.GetTempFolder(t, "bolt_store")
	defer cleanup()

	s := NewBoltDocStore(tmpPath)
	var err error
	_ = s
	_ = err
}
