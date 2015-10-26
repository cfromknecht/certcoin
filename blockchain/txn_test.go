package core

import (
	"testing"
)

func TestValid(t *testing.T) {
	sk := NewKey()
	txn := Txn{
		Body: TxnBody{
			From:      Address(sk.PublicKey),
			To:        "ANYONE",
			Value:     10000,
			PublicKey: sk.PublicKey,
		},
		Signature: CertcoinSignature{},
	}

	// Sign txn
	txn.Signature = Sign(txn.Body.Hash(), sk)

	if !txn.Valid() {
		t.Error("Excpected txn to be valid")
	}
}
