package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"
)

func NewPaymentTxn(from crypto.CertcoinPublicKey, to crypto.SHA256Sum, value uint64) Txn {
	return Txn{
		Type: Payment,
		Inputs: []Input{
			Input{
				PrevHash:  crypto.SHA256Sum{},
				PublicKey: from,
				Signature: crypto.CertcoinSignature{},
			},
		},
		Outputs: []Output{
			Output{
				Address: to,
				Value:   value,
			},
		},
	}
}

func (bc *Blockchain) ValidPaymentTxn(t Txn) bool {
	// Lookup amounts in UTXO
	return t.Type == Payment &&
		t.ValidNumInputs(0) &&
		t.ValidNumOutputs()
}
