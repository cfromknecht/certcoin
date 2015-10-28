package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"
	db "github.com/syndtr/goleveldb/leveldb"

	"log"
)

type SVP struct {
	HeaderDBPath string
	LastHeader   BlockHeader
}

func NewSVP() SVP {
	svp := SVP{
		HeaderDBPath: "db/header.db",
	}

	err := svp.WriteHeader(GenesisBlock().Header)
	if err != nil {
		log.Println(err)
		panic("Unable to add genesis block header to database")
	}

	return svp
}

func (s *SVP) ValidHeader(header BlockHeader) bool {
	if header.SeqNum == 0 {
		return header.PrevHash == crypto.SHA256Sum{} &&
			header.ValidPoW()
	}

	return s.LastHeader.Hash() == header.PrevHash &&
		header.ValidPoW()
}

func (s *SVP) WriteHeader(header BlockHeader) error {
	headerDB, err := db.OpenFile(s.HeaderDBPath, nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open header database")
	}
	defer headerDB.Close()

	headerJson := header.Json()
	hash := header.Hash()

	err = headerDB.Put(hash[:], headerJson, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Last header:", header)
	s.LastHeader = header

	return nil
}
