package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"strings"

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

	for i, b := range hs.Sum(nil)[:n] {
		c := hashAlphabet[b%alphabetSize]
		if i >= n {
			break
		}
		res[i] = c
	}
	return res
}

// GetUID returns new unique id (xid)
func GetUID() []byte {
	return xid.New().Bytes()
}

type SignedUIDGenerator struct {
	hasher  hash.Hash
	signLen int
	sep     string
}

func NewSignedUIDGenerator(secret []byte) *SignedUIDGenerator {
	return &SignedUIDGenerator{
		hasher:  hmac.New(sha256.New224, secret),
		signLen: 8,
		sep:     "_",
	}
}

func (s *SignedUIDGenerator) GetUID(n int) []byte {
	uid := Base58UID(n)
	sign := s.hasher.Sum(uid)
	signStr := base64.StdEncoding.EncodeToString(sign)[:s.signLen]
	return []byte(string(uid) + s.sep + signStr)

}

func (s *SignedUIDGenerator) Validate(data []byte) bool {
	dataParts := strings.SplitN(string(data), s.sep, 2)
	if len(dataParts) != 2 {
		return false
	}
	sign := s.hasher.Sum([]byte(dataParts[0]))
	signStr := base64.StdEncoding.EncodeToString(sign)[:s.signLen]
	return signStr == dataParts[1]
}
