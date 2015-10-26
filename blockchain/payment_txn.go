package core

import (
	"encoding/json"
	"log"
)

type PaymentTxn struct {
	TxnBase
	Body      PaymentTxnBody    `json:"body"`
	Signature CertcoinSignature `json:"signature"`
}

type PaymentTxnBody struct {
	From      string            `json:"from"`
	To        string            `json:"to"`
	Value     uint64            `json:"value"`
	PublicKey CertcoinPublicKey `json:"public_key"`
}

func (tb PaymentTxnBody) Hash() string {
	txnBodyJson, err := json.Marshal(tb)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal txn")
	}

	return CertcoinHash(txnBodyJson)
}

func NewPaymentTxn(from CertcoinPublicKey, to string, value uint64) PaymentTxn {
	return PaymentTxn{
		TxnBase: TxnBase{
			Type: Payment,
		},
		Body: PaymentTxnBody{
			From:      Address(from),
			To:        to,
			Value:     value,
			PublicKey: from,
		},
		Signature: CertcoinSignature{},
	}
}

func (t PaymentTxn) Valid() bool {
	return t.Type == Payment &&
		t.Body.From == Address(t.Body.PublicKey) &&
		len(t.Body.To) < 45 &&
		Verify(t.Body.Hash(), t.Signature, t.Body.PublicKey)
}

func (t PaymentTxn) TxnType() TxnType {
	return t.Type
}
