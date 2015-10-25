package blockchain

import (
	"crytpo/sha256"
	"io"
)

type SHA256Sum [32]byte

type Block struct {
	PrevHash SHA256Sum
}

func DoubleSHA256(s string) (err, SHA256Sum) {
	doubleSHA256([]bytes(s))
}

func DoubleSHA256(b []byte) (err, SHA256Sum) {
	h256 := sha256.New()
	h := h256.Sum(b)
	return h256.Sum(h)
}
