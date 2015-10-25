package main

import (
	"fmt"
	"github.com/cfromknecht/certcoin/blockchain"
)

func main() {
	fmt.Println("Hello world")
	fmt.Println("Hashing hello")
	fmt.Println(blockchain.CertHash64Str("hello"))
}
