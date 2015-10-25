package blockchain

import (
	"encoding/hex"
)

type Block struct {
	PrevHash SHA256Sum
}

func b64Encode(b []byte) string {
	return hex.EncodeToString(b)
}

func b64Decode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
