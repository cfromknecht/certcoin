package blockchain

import (
	"testing"
)

func TestCertHash(t *testing.T) {
	input := "hello"
	expectedOutput := "9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50"

	hash := CertHash64Str(input)
	if hash != expectedOutput {
		t.Error("Input should hash to " + expectedOutput + ", got: " + hash)
	}
}
