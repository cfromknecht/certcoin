package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"

	"log"
)

const (
	REGISTRATION_FEE = uint64(1000)
)

func NewRegisterTxn(onlineSecret, offlineSecret crypto.CertcoinSecretKey,
	source crypto.CertcoinPublicKey,
	identity Identity) Txn {

	fullName := identity.FullName()
	return Txn{
		Type: Register,
		Inputs: []Input{
			Input{
				PrevHash:  identity.Domain,
				PublicKey: onlineSecret.PublicKey,
				Signature: crypto.Sign(fullName, onlineSecret),
			},
			Input{
				PrevHash:  identity.Subdomain,
				PublicKey: offlineSecret.PublicKey,
				Signature: crypto.Sign(fullName, offlineSecret),
			},
			Input{
				PrevHash:  crypto.SHA256Sum{},
				PublicKey: source,
				Signature: crypto.CertcoinSignature{},
			},
		},
		Outputs: []Output{
			Output{
				Address: crypto.SHA256Sum{},
				Value:   REGISTRATION_FEE,
			},
		},
	}
}

func (bc *Blockchain) ValidRegisterTxn(t Txn) bool {
	if !t.ValidNumInputs(2) || !t.ValidNumOutputs() {
		return false
	}

	domain := t.Inputs[0].PrevHash.String()
	subdomain := t.Inputs[1].PrevHash.String()

	identity, err := NewIdentity(domain, subdomain)
	if err != nil {
		log.Println(err)
		return false
	}

	// Check that Identity is not registered
	// Check that Identity-PKs is not in accumulator

	return t.Type == Register &&
		t.Outputs[0].Value >= REGISTRATION_FEE &&
		crypto.Verify(identity.FullName(), t.Inputs[0].Signature, t.Inputs[0].PublicKey) &&
		crypto.Verify(identity.FullName(), t.Inputs[1].Signature, t.Inputs[1].PublicKey)
}
