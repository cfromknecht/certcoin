package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"
)

type GenerationTxn struct {
	TxnType TxnType
	Ins     []Input
	Outs    []Output
}

func NewGenerationTxn(to crypto.SHA256Sum) Txn {
	return Txn{
		Type:   Generation,
		Inputs: []Input{},
		Outputs: []Output{
			Output{
				Address: to,
				Value:   CURRENT_BLOCK_REWARD,
			},
		},
	}
}

func (bc *Blockchain) ValidGenerationTxn(t Txn) bool {
	return t.Type == Generation &&
		len(t.Inputs) == 0 &&
		t.ValidNumOutputs() &&
		t.Outputs[0].Value == CURRENT_BLOCK_REWARD
}
