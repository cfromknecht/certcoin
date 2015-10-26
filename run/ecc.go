package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"github.com/cfromknecht/certcoin/asm"
	"hash"
	"io"
	"math/big"
	"os"
)

func main2() {
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

func main() {
	numValues := (1 << 16) - 1
	acc := asm.NewAsyncAcc()
	witnesses := []asm.WitnessPath{}

	fmt.Print("[ADDING VALUES] |")
	for i := 0; i < numValues; i++ {
		wit := acc.Add(fmt.Sprintf("%d", i))
		witnesses = append(witnesses, wit)

		// Print progress bar
		divisor := numValues / 100
		if i%(10*divisor) == 0 && i != 0 {
			if i/divisor != 100 {
				fmt.Print(i / divisor)
			}
		} else if i%divisor == 0 {
			fmt.Print("=")
		}
	}

	for i, wit := range witnesses {
		fmt.Println("witness", i, wit)
	}

	fmt.Println("|\naccumulator for", numValues, "values:", acc)
}
