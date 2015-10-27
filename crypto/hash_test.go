package crypto

import (
	"testing"
)

func TestCertHash(t *testing.T) {
	input := "hello"
	expectedOutput := "lZXJ35AHUUjrBoYDZd8zWEt1v/eCpRDGzUiDpBmDPVA="

	hash := CertcoinHashStr(input)
	if hash != expectedOutput {
		t.Error("Input should hash to " + expectedOutput + ", got: " + hash)
	}
}
