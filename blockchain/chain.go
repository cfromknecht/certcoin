package blockchain

import (
	db "github.com/syndtr/goleveldb/leveldb"

	"errors"
	"log"
)

type Blockchain struct {
	SVP
	blockDB *db.DB
	utxnDB  *db.DB
}

func NewBlockchain() Blockchain {
	bc := Blockchain{
		SVP:     NewSVP(),
		blockDB: nil,
		utxnDB:  nil,
	}

	// Connect to block database
	blockConn, err := db.OpenFile("db/block.db", nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open block database")
	}
	defer func() { blockConn.Close() }()
	bc.blockDB = blockConn

	// Connect to utxn database
	utxnConn, err := db.OpenFile("db/utxn.db", nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open block database")
	}
	defer func() { utxnConn.Close() }()
	bc.utxnDB = utxnConn

	return bc
}

func (bc *Blockchain) VerifyAddBlock(b Block) error {
	if !bc.ValidBlock(b) {
		return errors.New("Cannot add invalid block")
	}

	err := bc.WriteHeader(b.Header)
	if err != nil {
		log.Println(err)
		return err
	}

	err = bc.WriteBlock(b)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (bc *Blockchain) WriteBlock(b Block) error {
	hash := b.Header.Hash()
	err := bc.blockDB.Put(hash[:], b.Json(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	batch := &db.Batch{}
	for _, txn := range b.Txns {
		txnHash := txn.Hash()
		batch.Put(txnHash[:], txn.Json())
	}
	err = bc.utxnDB.Write(batch, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (bc *Blockchain) ValidBlock(b Block) bool {
	if !b.ValidPoW() {
		return false
	}

	for i, txn := range b.Txns {
		if i == 0 && txn.Type != Generation {
			return false
		}

		switch txn.Type {
		case Generation:
			if !bc.ValidGenerationTxn(txn) {
				return false
			}
		case Payment:
			if !bc.ValidPaymentTxn(txn) {
				return false
			}
		case Register:
			if !bc.ValidRegisterTxn(txn) {
				return false
			}
		case Update:
			if !bc.ValidUpdateTxn(txn) {
				return false
			}
		case Revoke:
			if !bc.ValidRevokeTxn(txn) {
				return false
			}
		}
	}

	return true
}
