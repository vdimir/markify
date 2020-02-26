package util

import (
	"crypto/sha256"

	"github.com/rs/xid"
)

// In fact that we use `byte % alphabetSize` as index
// first 24 symbols (= 256 % 58), [a-y] have higher probapility than others.
// But it is ok for our purposes.
const hashAlphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"
const alphabetSize = byte(len(hashAlphabet))

// Base58UID random base58 identifier
func Base58UID(n int) []byte {
	hs := sha256.New224()

	hs.Write(GetUID())

	maxLen := hs.Size()
	if n > maxLen {
		n = maxLen
	}

	res := make([]byte, n)

	for i, b := range hs.Sum(nil) {
		c := hashAlphabet[b%alphabetSize]
		if i >= n {
			break
		}
		res[i] = c
	}
	return res
}

// GetUID retunst new unique id (xid)
func GetUID() []byte {
	return xid.New().Bytes()
}
