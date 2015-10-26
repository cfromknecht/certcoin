package core

type TxnType uint8

const (
	Generation TxnType = iota
	Payment
	Register
	Update
	Revoke
)

type Txn interface {
	TxnType() TxnType
	Valid() bool
}

type TxnBase struct {
	Type TxnType `json:"type"`
}
