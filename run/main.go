package main

import (
	"github.com/cfromknecht/certcoin/blockchain"

	"fmt"
)

func main() {

	b := blockchain.NewBlockchain()
	fmt.Println(b)
}
