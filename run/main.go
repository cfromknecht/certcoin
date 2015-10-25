package main

import (
	"fmt"
	"github.com/cfromknecht/certcoin"
)

func main() {
	fmt.Println("Hello world")
	fmt.Println("Hashing hello")
	fmt.Println(certcoin.Hash("hello"))
}
