package core

import (
	"encoding/json"
	"log"
)

const (
	UPDATE_FEE = uint64(100)
)

type UpdateTxn struct {
	TxnBase
	Body             UpdateTxnBody     `json:"body"`
	Signature        CertcoinSignature `json:"source_sig"`
	OfflineSignature CertcoinSignature `json:"offline_sig"`
}

type UpdateTxnBody struct {
	SourcePublicKey CertcoinPublicKey `json:"source_pk"`
	Value           uint64            `json:"value"`
	// PKI Informtion
	Identity        string            `json:"identity"`
	OnlinePublicKey CertcoinPublicKey `json:"online_pk"`
	OnlineSignature CertcoinSignature `json:"online_sig"`
}

func (t UpdateTxnBody) Hash() string {
	json, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal update txn body")
	}

	return CertcoinHash(json)
}

func NewUpdateTxn(onlineSecret CertcoinSecretKey,
	source CertcoinPublicKey,
	identity string) UpdateTxn {
	return UpdateTxn{
		TxnBase: TxnBase{
			Type: Update,
		},
		Body: UpdateTxnBody{
			SourcePublicKey: source,
			Value:           UPDATE_FEE,
			Identity:        identity,
			OnlinePublicKey: onlineSecret.PublicKey,
			OnlineSignature: Sign("", onlineSecret),
		},
		Signature:        CertcoinSignature{},
		OfflineSignature: CertcoinSignature{},
	}
}

func (t UpdateTxn) Valid() bool {
	// offlinePK := lookup from database
	return t.TxnType() == Update &&
		t.Body.Value >= UPDATE_FEE &&
		Verify(t.Body.Hash(), t.Signature, t.Body.SourcePublicKey) &&
		Verify("", t.Body.OnlineSignature, t.Body.OnlinePublicKey) &&
		true
	//Verify(t.Body.Hash(), t.OfflineSignature, offlinePK)
}

func (t UpdateTxn) TxnType() TxnType {
	return t.Type
}
