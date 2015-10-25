package blockchain

import (
	"crypto/sha256"
	"encoding/base64"
)

type SHA256Sum []byte

type Block struct {
	PrevHash SHA256Sum
}

func CertHash(b []byte) SHA256Sum {
	h256 := sha256.New()
	h := h256.Sum(b)
	return h256.Sum(h)
}

func CertHash64(b []byte) SHA256Sum {
	h := CertHash(b)
	return b64Encode(h)
}

func CertHashStr(s string) SHA256Sum {
	return CertHash([]byte(s))
}

func CertHash64Str(s string) SHA256Sum {
	h := CertHashStr(s)
	return b64Encode(h)
}

func b64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func b64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeStr(s)
}
