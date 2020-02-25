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

// func Base58(raw []byte) []byte {
// 	res := make([]byte, 3*len(raw))
// 	for _, b := range raw {
// 		c := hashAlphabet[b%alphabetSize]
// 	}
// }

// Base58UID hash string and encode with base58
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
