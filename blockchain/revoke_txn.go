package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"

	"fmt"
	"log"
)

const (
	REVOKE_FEE = uint64(100)
)

func NewRevokeTxn(onlineSecret, offlineSecret crypto.CertcoinSecretKey,
	source crypto.CertcoinPublicKey,
	identity Identity) Txn {

	sig := crypto.Sign(identity.FullName(), onlineSecret)
	log.Println("Signature:", sig)
	verifies := crypto.Verify(identity.FullName(), sig, onlineSecret.PublicKey)
	log.Println("PublicKey:", onlineSecret.PublicKey)
	log.Println("Signature:", sig)
	log.Println("Verifies:", verifies)

	return Txn{
		Type: Revoke,
		Inputs: []Input{
			Input{
				PrevHash:  identity.Domain,
				PublicKey: onlineSecret.PublicKey,
				Signature: sig,
			},
			Input{
				PrevHash:  identity.Subdomain,
				PublicKey: offlineSecret.PublicKey,
				Signature: sig,
			},
			Input{
				PrevHash:  crypto.SHA256Sum{},
				PublicKey: source,
				Signature: sig,
			},
		},
		Outputs: []Output{
			Output{
				Address: crypto.SHA256Sum{},
				Value:   REVOKE_FEE,
			},
		},
	}
}

func (bc *Blockchain) ValidRevokeTxn(t Txn) bool {
	if !t.ValidNumInputs(2) || !t.ValidNumOutputs() {
		return false
	}

	fmt.Println("Txn: %v", t)

	domain := t.Inputs[0].PrevHash.String()
	subdomain := t.Inputs[1].PrevHash.String()

	identity, err := NewIdentity(domain, subdomain)
	if err != nil {
		return false
	}

	// offlinePK := lookup from database

	return t.Type == Revoke &&
		t.Outputs[0].Value >= REVOKE_FEE &&
		crypto.Verify(identity.FullName(), t.Inputs[0].Signature, t.Inputs[0].PublicKey)
}
