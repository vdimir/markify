package util

import (
	"crypto/sha256"
	"encoding/binary"
	"time"

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

// BaseHashEncode hash string and encode with base58
func BaseHashEncode(raw []byte, n int) ([]byte, int64) {
	hs := sha256.New224()

	hs.Write(raw)
	seed := time.Now().UnixNano()
	binary.Write(hs, binary.LittleEndian, seed)

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
	return res, seed
}

// GetUID retunst new unique id (xid)
func GetUID() []byte {
	return xid.New().Bytes()
}
