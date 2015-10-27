package blockchain

import (
	db "github.com/syndtr/goleveldb/leveldb"

	"log"
)

type SVP struct {
	headerDB   *db.DB
	LastHeader BlockHeader
}

func NewSVP() SVP {
	svp := SVP{
		headerDB:   nil,
		LastHeader: GenesisBlock().Header,
	}

	headerConn, err := db.OpenFile("db/header.db", nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open header database")
	}
	defer func() { headerConn.Close() }()
	svp.headerDB = headerConn

	return svp
}

func (s *SVP) WriteHeader(header BlockHeader) error {
	headerJson := header.Json()
	hash := header.Hash()

	err := s.headerDB.Put(hash[:], headerJson, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	s.LastHeader = header

	return nil
}
