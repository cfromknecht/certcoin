package core

import (
	"encoding/json"
	"log"
)

const (
	REVOKE_FEE = uint64(100)
)

type RevokeTxn struct {
	TxnBase
	Body             RevokeTxnBody     `json:"body"`
	Signature        CertcoinSignature `json:"source_sig"`
	OfflineSignature CertcoinSignature `json:"offline_sig"`
}

type RevokeTxnBody struct {
	SourcePublicKey CertcoinPublicKey `json:"source_pk"`
	Value           uint64            `json:"value"`
	// PKI Informtion
	Identity string `json:"identity"`
}

func (t RevokeTxnBody) Hash() string {
	json, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal revoke txn body")
	}

	return CertcoinHash(json)
}

func NewRevokeTxn(offlineSecret CertcoinSecretKey,
	source CertcoinPublicKey,
	identity string) RevokeTxn {
	return RevokeTxn{
		TxnBase: TxnBase{
			Type: Revoke,
		},
		Body: RevokeTxnBody{
			SourcePublicKey: source,
			Value:           REVOKE_FEE,
			Identity:        identity,
		},
		Signature:        CertcoinSignature{},
		OfflineSignature: CertcoinSignature{},
	}
}

func (t RevokeTxn) Valid() bool {
	// offlinePK := lookup from database
	return t.TxnType() == Revoke &&
		t.Body.Value >= REVOKE_FEE &&
		Verify(t.Body.Hash(), t.Signature, t.Body.SourcePublicKey) &&
		true
	//Verify(t.Body.Hash(), t.OfflineSignature, offlinePK)
}

func (t RevokeTxn) TxnType() TxnType {
	return t.Type
}
