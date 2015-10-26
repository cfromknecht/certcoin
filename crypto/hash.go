package core

import (
	"crypto/sha256"
	"encoding/base64"
)

func CertcoinHash(b []byte) string {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])

	return b64Encode(hash2[:])
}

func CertcoinHashStr(s string) string {
	return CertcoinHash([]byte(s))
}

func b64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func b64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
