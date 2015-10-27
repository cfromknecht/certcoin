package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"
)

const (
	REGISTRATION_FEE = uint64(1000)
)

func NewRegisterTxn(onlineSecret, offlineSecret crypto.CertcoinSecretKey,
	source crypto.CertcoinPublicKey,
	identity Identity) Txn {

	fullName := identity.SubdomainStr() + identity.DomainStr()
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
	if len(t.Inputs) < 3 || !(len(t.Outputs) == 1 && len(t.Outputs) == 2) {
		return false
	}

	identity, err := NewIdentity(string(t.Inputs[0].PrevHash[:]), string(t.Inputs[1].PrevHash[:]))
	if err != nil {
		return false
	}

	// Check that Identity is not registered
	// Check that Identity-PKs is not in accumulator

	return t.Type == Register &&
		t.Outputs[0].Value >= REGISTRATION_FEE &&
		crypto.Verify(identity.FullNameStr(), t.Inputs[0].Signature, t.Inputs[0].PublicKey) &&
		crypto.Verify(identity.FullNameStr(), t.Inputs[1].Signature, t.Inputs[1].PublicKey)
}
