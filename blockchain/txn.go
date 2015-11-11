package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"

	"encoding/json"
	"log"
)

type TxnType uint8

const (
	Generation TxnType = iota
	Payment
	Register
	Update
	Revoke
)

type Txn struct {
	Type    TxnType  `json:"txn_type"`
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

type Input struct {
	PrevHash  crypto.SHA256Sum         `json:"hash"`
	PublicKey crypto.CertcoinPublicKey `json:"public"`
	Signature crypto.CertcoinSignature `json:"sig"`
}

type Output struct {
	Address crypto.SHA256Sum `json:"address"`
	Value   uint64           `json:"value"`
}

func (t Txn) Json() []byte {
	txnJson, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal txn")
	}
	return txnJson
}

func (t Txn) Hash() crypto.SHA256Sum {
	return crypto.CertcoinHash(t.Json())
}

func (t Txn) ValidNumInputs(reserved int) bool {
	return len(t.Inputs) >= reserved+1
}

func (t Txn) ValidNumOutputs() bool {
	return len(t.Outputs) == 1 ||
		len(t.Outputs) == 2
}
