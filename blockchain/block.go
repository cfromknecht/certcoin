package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

type SHA256Sum [32]byte

type Block struct {
	PrevHash SHA256Sum
}

func CertHash(b []byte) SHA256Sum {
	hash1 := sha256.Sum256(b)
	return sha256.Sum256(hash1[:])
}

func CertHash64(b []byte) string {
	h := CertHash(b)
	return b64Encode(h[:])
}

func CertHashStr(s string) SHA256Sum {
	return CertHash([]byte(s))
}

func CertHash64Str(s string) string {
	h := CertHashStr(s)
	return b64Encode(h[:])
}

func b64Encode(b []byte) string {
	return hex.EncodeToString(b)
}

func b64Decode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
