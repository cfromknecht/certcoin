package blockchain

import (
	db "github.com/syndtr/goleveldb/leveldb"

	"log"
)

type Blockchain struct {
	SVP
	BlockDBPath string
	UTxnDBPath  string
}

func NewBlockchain() Blockchain {
	bc := Blockchain{
		SVP:         NewSVP(),
		BlockDBPath: "db/block.db",
		UTxnDBPath:  "db/utxn.db",
	}

	g := GenesisBlock()
	for !g.Header.ValidPoW() {
		g.Header.Nonce += 1
	}

	if bc.ValidBlock(g) {
		err := bc.WriteBlock(g)
		if err != nil {
			log.Println(err)
			panic("Unable to add genesis block to database")
		}
	} else {
		panic("Genesis block invalid")
	}

	return bc
}

func (bc *Blockchain) ValidBlock(b Block) bool {
	return bc.ValidHeader(b.Header) &&
		bc.ValidTxns(b)
}

func (bc *Blockchain) WriteBlock(b Block) error {
	err := bc.WriteHeader(b.Header)
	if err != nil {
		log.Println(err)
		return err
	}

	blockDB, err := db.OpenFile(bc.BlockDBPath, nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open block database")
	}
	defer blockDB.Close()

	hash := b.Header.Hash()
	err = blockDB.Put(hash[:], b.Json(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	utxnDB, err := db.OpenFile(bc.UTxnDBPath, nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open utxn database")
	}
	defer utxnDB.Close()

	batch := &db.Batch{}
	for _, txn := range b.Txns {
		txnHash := txn.Hash()
		batch.Put(txnHash[:], txn.Json())
	}
	err = utxnDB.Write(batch, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (bc *Blockchain) ValidTxns(b Block) bool {
	// Check validity of each transaction
	for i, txn := range b.Txns {
		if i == 0 && txn.Type != Generation {
			log.Println("First is not a GenerationTxn")
			return false
		}

		switch txn.Type {
		case Generation:
			if !bc.ValidGenerationTxn(txn) {
				log.Println("Invalid GenerationTxn")
				return false
			}
		case Payment:
			if !bc.ValidPaymentTxn(txn) {
				log.Println("Invalid PaymentTxn")
				return false
			}
		case Register:
			if !bc.ValidRegisterTxn(txn) {
				log.Println("Invalid RegisterTxn")
				return false
			}
		case Update:
			if !bc.ValidUpdateTxn(txn) {
				log.Println("Invalid UpdateTxn")
				return false
			}
		case Revoke:
			if !bc.ValidRevokeTxn(txn) {
				log.Println("Invalid RevokeTxn")
				return false
			}
		}
	}

	// Check that accumulator is constructed properly
	prevAcc := bc.SVP.LastHeader.PKIAcc
	for _, txn := range b.Txns {
		switch txn.Type {
		case Register:
			domain := txn.Inputs[0].PrevHash.String()
			subdomain := txn.Inputs[1].PrevHash.String()
			identity, err := NewIdentity(domain, subdomain)
			if err != nil {
				log.Println(err)
				return false
			}
			fullName := identity.FullName()

			y := fullName +
				string(txn.Inputs[0].PublicKey.X.Bytes()) +
				string(txn.Inputs[0].PublicKey.Y.Bytes())
			prevAcc.Add(y)
		default:
			continue
		}
	}

	return compareSlice(prevAcc, b.Header.PKIAcc)
}

func compareSlice(s1, s2 []string) bool {
	if s1 == nil && s2 == nil {
		return true
	}

	if s1 == nil || s2 == nil {
		return false
	}

	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}

	}

	return true
}
