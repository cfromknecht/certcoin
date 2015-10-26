package core

import (
	"encoding/json"
	"log"
)

const (
	REGISTRATION_FEE = uint64(1000)
)

type RegisterTxn struct {
	TxnBase
	Body      RegisterTxnBody
	Signature CertcoinSignature
}
type RegisterTxnBody struct {
	SourcePublicKey CertcoinPublicKey `json:"source_pk"`
	Value           uint64            `json:"value"`
	// PKI Informtion
	Identity         string            `json:"identity"`
	OnlinePublicKey  CertcoinPublicKey `json:"online_pk"`
	OfflinePublicKey CertcoinPublicKey `json:"offline_pk"`
	OnlineSignature  CertcoinSignature `json:"online_sig"`
	OfflineSignature CertcoinSignature `json:"offline_sig"`
}

func (b RegisterTxnBody) Hash() string {
	json, err := json.Marshal(b)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal registration txn")
	}

	return CertcoinHash(json)
}

func NewRegisterTxn(onlineSecret, offlineSecret CertcoinSecretKey,
	source CertcoinPublicKey,
	identity string) RegisterTxn {
	return RegisterTxn{
		TxnBase: TxnBase{
			Type: Register,
		},
		Body: RegisterTxnBody{
			SourcePublicKey:  source,
			Value:            REGISTRATION_FEE,
			Identity:         identity,
			OnlinePublicKey:  onlineSecret.PublicKey,
			OfflinePublicKey: offlineSecret.PublicKey,
			OnlineSignature:  Sign("", onlineSecret),
			OfflineSignature: Sign("", offlineSecret),
		},
		Signature: CertcoinSignature{},
	}
}

func (r RegisterTxn) Valid() bool {
	return r.TxnType() == Register &&
		r.Body.Value >= REGISTRATION_FEE &&
		Verify(r.Body.Hash(), r.Signature, r.Body.SourcePublicKey) &&
		Verify("", r.Body.OnlineSignature, r.Body.OnlinePublicKey) &&
		Verify("", r.Body.OfflineSignature, r.Body.OfflinePublicKey)
}

func (r RegisterTxn) TxnType() TxnType {
	return r.Type
}
