package main

import (
	"accumulator/certcoin"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"hash"
	"io"
	"math/big"
	"os"
)

func main() {
	curve := elliptic.P256()

	privateKey := new(ecdsa.PrivateKey)
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var publicKey ecdsa.PublicKey
	publicKey = privateKey.PublicKey

	fmt.Println("Private Key :")
	fmt.Printf("%x \n", privateKey)

	fmt.Println("Public Key :")
	fmt.Printf("%x \n", publicKey)

	var h hash.Hash
	h = md5.New()
	r := big.NewInt(0)
	s := big.NewInt(0)

	io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	signhash := h.Sum(nil)

	r, s, serr := ecdsa.Sign(rand.Reader, privateKey, signhash)
	if serr != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)

	fmt.Printf("Signature : %x\n", signature)

	verifyStatus := ecdsa.Verify(&publicKey, signhash, r, s)
	fmt.Println(verifyStatus)
}
