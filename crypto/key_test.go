package crypto

import (
	"testing"
)

func TestNewKey(t *testing.T) {
	keyPair := NewKey()

	msg := "hello"
	sig := Sign(msg, keyPair)
	verifies := Verify(msg, sig, keyPair.PublicKey)

	if !verifies {
		t.Error("Signature for key pair does not verify")
	}
}
