package main

import (
	"fmt"
	"github.com/cfromknecht/certcoin/core"
)

func main() {

	keypair := core.NewKey()
	sig := core.Sign("hello", keypair)
	verifies := core.Verify("hello", sig, keypair.PublicKey)

	fmt.Println("Signature verifies:", verifies)

	g := core.GenesisBlock()
	fmt.Println(fmt.Sprintf("Genesis Block:\n%v", g.Json()))
	fmt.Println(g.Hash())

	b1 := core.NewBlock(g, core.Address(keypair.PublicKey))
	fmt.Println(fmt.Sprintf("Block 1:\n%v", b1.Json()))
	fmt.Println(b1.Hash())

	core.Mine()
}
