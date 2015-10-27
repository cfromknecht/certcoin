package blockchain

import (
	"github.com/cfromknecht/certcoin/asm"
	"github.com/cfromknecht/certcoin/crypto"

	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	HASH_GENESIS_BLOCK = "hwm5tjMI9ZsvG2VwNKeMtJq7PDLDvPM/hLFQ2VYdE88="
)

var (
	CURRENT_DIFFICULTY   = uint64(16)
	CURRENT_BLOCK_REWARD = uint64(50000000)
)

type BlockHeader struct {
	SeqNum     uint64       `json:"seq_num"`
	PrevHash   string       `json:"prev_hash"`
	TxnAcc     asm.AsyncAcc `json:"txn_acc"`
	PKIAcc     asm.AsyncAcc `json:"pki_acc"`
	Time       time.Time    `json:"time"`
	Difficulty uint64       `json:"difficulty"`
	Nonce      uint64       `json:"nonce:"`
}

type Block struct {
	Header BlockHeader `json:"header"`
	Txns   []Txn       `json:"txns"`
}

func (b Block) Json() []byte {
	blockJson, err := json.Marshal(b)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal block")
	}

	return blockJson
}

func NewBlock(prev Block, minerAddress crypto.SHA256Sum) Block {
	return Block{
		Header: BlockHeader{
			SeqNum:     prev.Header.SeqNum + 1,
			PrevHash:   prev.Header.Hash().String(),
			TxnAcc:     prev.Header.TxnAcc,
			PKIAcc:     prev.Header.PKIAcc,
			Time:       time.Now(),
			Difficulty: CURRENT_DIFFICULTY,
			Nonce:      0,
		},
		Txns: []Txn{
			NewGenerationTxn(minerAddress),
		},
	}
}

func GenesisBlock() Block {
	return Block{
		Header: BlockHeader{
			SeqNum:     1,
			PrevHash:   "",
			TxnAcc:     asm.NewAsyncAcc(),
			PKIAcc:     asm.NewAsyncAcc(),
			Time:       time.Now(),
			Difficulty: CURRENT_DIFFICULTY,
			Nonce:      0,
		},
		Txns: []Txn{},
	}
}

func (b BlockHeader) Json() []byte {
	headerJson, err := json.Marshal(b)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal block header")
	}

	return headerJson
}

func (b BlockHeader) Hash() crypto.SHA256Sum {
	return crypto.CertcoinHash(b.Json())
}

func (b Block) ValidPoW() bool {
	difficulty := b.Header.Difficulty
	zeroBytes := difficulty / 8
	bitOffset := difficulty % 8

	h, err := crypto.B64Decode(b.Header.Hash().String())
	if err != nil {
		log.Println(err)
		panic("Unable to base64 decode hash")
	}

	for i, c := range h {
		if uint64(i) < zeroBytes {
			if c != 0 {
				return false
			}
		} else {
			break
		}
	}

	c := h[zeroBytes]
	for j := 0; j < 8; j++ {
		if uint64(j) < bitOffset {
			if (c >> (7 - uint64(j)) & 1) != 0 {
				return false
			}
		} else {
			break
		}
	}

	return true
}

func Mine() {
	prev := GenesisBlock()
	fromKey := crypto.NewKey()
	toKey := crypto.NewKey()

	online := crypto.NewKey()
	offline := crypto.NewKey()

	//for {
	b := NewBlock(prev, crypto.Address(fromKey.PublicKey))

	// Create and sign txn
	txn := NewPaymentTxn(fromKey.PublicKey, crypto.Address(toKey.PublicKey), 10)
	txn.Inputs[0].Signature = crypto.Sign("", fromKey)
	b.Txns = append(b.Txns, txn)

	identity, err := NewIdentity("certcoin.net", "")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Create and sign registration txn
	rtxn := NewRegisterTxn(online, offline, fromKey.PublicKey, identity)
	rtxn.Inputs[2].Signature = crypto.Sign("", fromKey)
	b.Txns = append(b.Txns, rtxn)

	newOnline := crypto.NewKey()
	// Create and sign update txn
	utxn := NewUpdateTxn(newOnline, offline, fromKey.PublicKey, identity)
	utxn.Inputs[1].Signature = crypto.Sign("", offline)
	utxn.Inputs[2].Signature = crypto.Sign("", fromKey)
	b.Txns = append(b.Txns, utxn)

	// Create and sign revoke txn
	vtxn := NewRevokeTxn(newOnline, offline, fromKey.PublicKey, identity)
	vtxn.Inputs[0].Signature = crypto.Sign("", offline)
	vtxn.Inputs[1].Signature = crypto.Sign("", fromKey)
	b.Txns = append(b.Txns, vtxn)

	for !b.ValidPoW() {
		b.Header.Nonce += 1
	}

	fmt.Println(fmt.Sprintf("%v", b.Json()))
	fmt.Println(b.Header.Hash())
	prev = b
	//}
}
