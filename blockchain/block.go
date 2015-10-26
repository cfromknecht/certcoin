package core

import (
	"github.com/cfromknecht/certcoin/asm"

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
	SeqNum      uint64       `json:"seq_num"`
	PrevHash    string       `json:"prev_hash"`
	MerkelRoot  asm.AsyncAcc `json:"merkle_root"`
	Accumulator asm.AsyncAcc `json:"accumulator"`
	Time        time.Time    `json:"time"`
	Difficulty  uint64       `json:"difficulty"`
	Nonce       uint64       `json:"nonce:"`
}

type Block struct {
	Header       BlockHeader   `json:"header"`
	GenTxn       GenerationTxn `json:"gen_txn"`
	PaymentTxns  []PaymentTxn  `json:"payment_txns"`
	RevokeTxns   []RevokeTxn   `json:"revoke_txns"`
	UpdateTxns   []UpdateTxn   `json:"update_txns"`
	RegisterTxns []RegisterTxn `json:"register_txns"`
}

func (b Block) Json() string {
	json, err := json.Marshal(b)
	if err != nil {
		log.Println(err)
		panic("Unable to json.Marshal block")
	}

	return string(json)
}

func NewBlock(prev Block, minerAddress string) Block {
	return Block{
		Header: BlockHeader{
			SeqNum:      prev.Header.SeqNum + 1,
			PrevHash:    prev.Hash(),
			MerkelRoot:  prev.Header.MerkelRoot,
			Accumulator: prev.Header.Accumulator,
			Time:        time.Now(),
			Difficulty:  CURRENT_DIFFICULTY,
			Nonce:       0,
		},
		Txns: []Txn{
			NewGenerationTxn(minerAddress),
		},
	}
}

func GenesisBlock() Block {
	return Block{
		Header: BlockHeader{
			SeqNum:      1,
			PrevHash:    "",
			MerkelRoot:  asm.NewAsyncAcc(),
			Accumulator: asm.NewAsyncAcc(),
			Time:        time.Now(),
			Difficulty:  CURRENT_DIFFICULTY,
			Nonce:       0,
		},
		Txns: []Txn{},
	}
}

func (b Block) Hash() string {
	headerJson, err := json.Marshal(b.Header)
	if err != nil {
		log.Println(err)
		panic("Unable to marshal block")
	}

	return CertcoinHash(headerJson)
}

func (b Block) Valid() bool {
	if !b.ValidPoW() {
		fmt.Println("Invalid PoW")
		return false
	}

	if !b.ValidTxns() {
		fmt.Println("Invalid txn")
		return false
	}

	return true
}

func (b Block) ValidPoW() bool {
	difficulty := b.Header.Difficulty
	zeroBytes := difficulty / 8
	bitOffset := difficulty % 8

	h, err := b64Decode(b.Hash())
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
	fromKey := NewKey()
	toKey := NewKey()

	online := NewKey()
	offline := NewKey()

	for {
		b := NewBlock(prev, Address(fromKey.PublicKey))

		// Create and sign txn
		txn := NewPaymentTxn(fromKey.PublicKey, Address(toKey.PublicKey), 10)
		txn.Signature = Sign(txn.Body.Hash(), fromKey)
		b.Txns = append(b.Txns, txn)

		// Create and sign registration txn
		rtxn := NewRegisterTxn(online, offline, fromKey.PublicKey, "certcoin.net")
		rtxn.Signature = Sign(rtxn.Body.Hash(), fromKey)
		b.Txns = append(b.Txns, rtxn)

		newOnline := NewKey()
		// Create and sign update txn
		utxn := NewUpdateTxn(newOnline, fromKey.PublicKey, "certcoin.net")
		utxn.Signature = Sign(utxn.Body.Hash(), fromKey)
		utxn.OfflineSignature = Sign(utxn.Body.Hash(), offline)
		b.Txns = append(b.Txns, utxn)

		// Create and sign revoke txn
		vtxn := NewRevokeTxn(offline, fromKey.PublicKey, "certcoin.net")
		vtxn.Signature = Sign(vtxn.Body.Hash(), fromKey)
		vtxn.OfflineSignature = Sign(vtxn.Body.Hash(), offline)
		b.Txns = append(b.Txns, vtxn)

		for !b.ValidPoW() {
			b.Header.Nonce += 1
		}

		fmt.Println(fmt.Sprintf("%v", b.Json()))
		fmt.Println(b.Hash())
		fmt.Println("Valid?:", b.Valid())
		prev = b
	}
}

func (b Block) ValidTxns() bool {
	for i, txn := range b.Txns {
		if i == 0 {
			if txn.TxnType() != Generation || !txn.Valid() {
				return false
			}
		} else if !txn.Valid() {
			return false
		}
	}

	return true
}
