package core

type GenerationTxn struct {
	TxnBase
	To    string
	Value uint64
}

func NewGenerationTxn(to string) GenerationTxn {
	return GenerationTxn{
		TxnBase: TxnBase{
			Type: Generation,
		},
		To:    to,
		Value: CURRENT_BLOCK_REWARD,
	}
}

func (t GenerationTxn) Valid() bool {
	return t.TxnType() == Generation &&
		len(t.To) < 45 &&
		t.Value == CURRENT_BLOCK_REWARD
}

func (t GenerationTxn) TxnType() TxnType {
	return t.Type
}
