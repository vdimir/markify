package kvstore

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vdimir/markify/testutil"
	bolt "go.etcd.io/bbolt"
)

func closeCheck(t *testing.T, s Store) {
	assert.NoError(t, s.Close())
}

func TestBoltStore(t *testing.T) {
	tmpPath, cleanup := testutil.GetTempFolder(t, "bolt_store")
	defer cleanup()

	boltPath := path.Join(tmpPath, "bolt.db")
	keyStore, err := NewBoltStorage(boltPath, bolt.Options{})
	defer closeCheck(t, keyStore)

	require.NoError(t, err)

	err = keyStore.Save([]byte("abc"), []byte("123"))
	assert.NoError(t, err)

	data, err := keyStore.Load([]byte("xyz"))
	assert.NoError(t, err)
	assert.Nil(t, data)
	ts0, err := keyStore.Timestamp([]byte("xyz"))
	assert.NoError(t, err)
	assert.Zero(t, ts0)

	data, err = keyStore.Load([]byte("abc"))
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("123"))
	ts1, err := keyStore.Timestamp([]byte("abc"))
	assert.NoError(t, err)
	assert.NotZero(t, ts1)

	err = keyStore.Save([]byte("abc"), []byte("456"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "already exists")

	err = keyStore.Update([]byte("abc"), []byte("456"))
	require.NoError(t, err)

	data, err = keyStore.Load([]byte("abc"))
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("456"))
	ts2, err := keyStore.Timestamp([]byte("abc"))
	assert.NoError(t, err)
	assert.NotZero(t, ts2)
	assert.Greater(t, ts2, ts1)
}

func TestBoltNewEntry(t *testing.T) {
	tmpPath, cleanup := testutil.GetTempFolder(t, "bolt_store")
	defer cleanup()

	boltPath := path.Join(tmpPath, "bolt.db")
	keyStore, err := NewBoltStorage(boltPath, bolt.Options{})
	require.NoError(t, err)
	defer closeCheck(t, keyStore)

	k1, err := keyStore.NewEntry([]byte("abc"))
	assert.NoError(t, err)
	k2, err := keyStore.NewEntry([]byte("bca"))
	assert.NoError(t, err)
	k3, err := keyStore.NewEntry([]byte("abc"))
	assert.NoError(t, err)

	var data []byte
	data, err = keyStore.Load(k1)
	assert.Equal(t, data, []byte("abc"))
	assert.NoError(t, err)
	data, err = keyStore.Load(k2)
	assert.Equal(t, data, []byte("bca"))
	assert.NoError(t, err)
	data, err = keyStore.Load(k3)
	assert.Equal(t, data, []byte("abc"))
	assert.NoError(t, err)

	assert.NotEqual(t, k1, k2)
	assert.NotEqual(t, k1, k3)
	assert.NotEqual(t, k2, k3)

	n := 100
	keys := map[string]bool{}
	for i := 0; i < n; i++ {
		k, err := keyStore.NewEntry([]byte("abc"))
		assert.NoError(t, err)
		keys[string(k)] = true
		assert.GreaterOrEqual(t, len(k), 12)
	}
	assert.Len(t, keys, n)
}

func BenchmarkBoltNewEntry(b *testing.B) {
	tmpPath, cleanup := testutil.GetTempFolder(b, "bolt_store")
	defer cleanup()

	boltPath := path.Join(tmpPath, "bolt.db")
	keyStore, err := NewBoltStorage(boltPath, bolt.Options{})
	require.NoError(b, err)

	keys := map[string]bool{}
	for i := 0; i < b.N; i++ {
		k, err := keyStore.NewEntry([]byte("abc"))
		assert.NoError(b, err)
		keys[string(k)] = true
		assert.GreaterOrEqual(b, len(k), 12)
	}
	assert.Len(b, keys, b.N)
}
