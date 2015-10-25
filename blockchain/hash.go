package blockchain

import (
	"crypto/sha256"
)

type SHA256Sum [32]byte

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
