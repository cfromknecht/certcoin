package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
)

var (
	CURVE = elliptic.P256()
)

type CertcoinPublicKey struct {
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

type CertcoinSecretKey struct {
	PublicKey CertcoinPublicKey `json:"public_key"`
	D         *big.Int          `json:"d"`
}

type CertcoinSignature struct {
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func NewKey() CertcoinSecretKey {
	keyPair := &ecdsa.PrivateKey{}
	keyPair, err := ecdsa.GenerateKey(CURVE, rand.Reader)
	if err != nil {
		log.Println(err)
		panic("Failed to created CertcoinKeyPair")
	}

	secret := CertcoinSecretKey{
		PublicKey: CertcoinPublicKey{
			X: keyPair.PublicKey.X,
			Y: keyPair.PublicKey.Y,
		},
		D: keyPair.D,
	}

	return secret
}

func Sign(m string, sk CertcoinSecretKey) CertcoinSignature {
	secret := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: CURVE,
			X:     sk.PublicKey.X,
			Y:     sk.PublicKey.Y,
		},
		D: sk.D,
	}

	h := CertcoinHashStr(m)
	r, s, err := ecdsa.Sign(rand.Reader, secret, h[:])
	if err != nil {
		log.Println(err)
		panic("Failed to sign message")
	}

	return CertcoinSignature{
		R: r,
		S: s,
	}
}

func Verify(m string, sig CertcoinSignature, pk CertcoinPublicKey) bool {
	public := ecdsa.PublicKey{
		Curve: CURVE,
		X:     pk.X,
		Y:     pk.Y,
	}

	h := CertcoinHashStr(m)
	return ecdsa.Verify(&public, h[:], sig.R, sig.S)
}

func Address(pk CertcoinPublicKey) SHA256Sum {
	s, err := json.Marshal(pk)
	if err != nil {
		log.Println(err)
		panic("Failed to marshal public key")
	}

	return CertcoinHash(s)
}
