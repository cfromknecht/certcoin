package crypto

import (
	"crypto/sha256"
	"encoding/base64"
)

type SHA256Sum [32]byte

func (s SHA256Sum) String() string {
	return B64Encode(s[:])
}

func CertcoinHash(b []byte) SHA256Sum {
	hash1 := sha256.Sum256(b)
	hash2 := sha256.Sum256(hash1[:])

	return hash2
}

func CertcoinHashStr(s string) SHA256Sum {
	return CertcoinHash([]byte(s))
}

func B64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func B64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
