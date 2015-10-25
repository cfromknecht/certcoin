package certcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"math/big"
	"os"
	"strings"
)

type RegistrationMessage struct {
	RequestedDomainName string

	OnlinePK       ecdsa.PublicKey
	OnlinePKSigned string

	OfflinePK       ecdsa.PublicKey
	OfflinePKSigned string
}

func signPublicKey(keyPair *ecdsa.PrivateKey) string {
	var h hash.Hash
	h = sha256.New()
	r := big.NewInt(0)
	s := big.NewInt(0)

	// hash public key before signing
	pkString := fmt.Sprintf("%x", keyPair.PublicKey)
	io.WriteString(h, pkString)
	pkHash := h.Sum(nil)

	// sign public key
	r, s, err := ecdsa.Sign(rand.Reader, keyPair, pkHash)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// format hex strings like this: "<r> <s>"
	pkSignatureNew := fmt.Sprintf("%x", r.Bytes()) + " " + fmt.Sprintf("%x", s.Bytes())
	return pkSignatureNew
}

func verifySignedPublicKey(pkSignature string, publicKey ecdsa.PublicKey) bool {
	rs := strings.Split(pkSignature, " ")
	rString := rs[0]
	sString := rs[1]

	r, success := new(big.Int).SetString(rString, 16)
	if !success {
		fmt.Println("Failed to set big.Int for with string")
		os.Exit(1)
	}
	s, success := new(big.Int).SetString(sString, 16)
	if !success {
		fmt.Println("Failed to set big.Int for with string")
		os.Exit(1)
	}

	h := sha256.New()
	pkString := fmt.Sprintf("%x", publicKey)
	io.WriteString(h, pkString)
	pkHash := h.Sum(nil)

	verifyStatus := ecdsa.Verify(&publicKey, pkHash, r, s)

	return verifyStatus
}

func main_test() {
	curve := elliptic.P256()

	onlineKey := new(ecdsa.PrivateKey)
	onlineKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	offlineKey := new(ecdsa.PrivateKey)
	offlineKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Online Key:", onlineKey.PublicKey)
	fmt.Println("Offline Key:", offlineKey.PublicKey)

	onlinePKSigned := signPublicKey(onlineKey)
	offlinePKSigned := signPublicKey(offlineKey)

	fmt.Println("Online Key Signed:", onlinePKSigned)
	fmt.Println("Offline Key Signed:", offlinePKSigned)

	onlinePKVerified := verifySignedPublicKey(onlinePKSigned, onlineKey.PublicKey)
	offlinePKVerified := verifySignedPublicKey(offlinePKSigned, offlineKey.PublicKey)

	fmt.Println("Online Key Verify:", onlinePKVerified)
	fmt.Println("Offline Key Verify:", offlinePKVerified)
}
