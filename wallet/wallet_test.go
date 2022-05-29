package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey = "30770201010420df9ea67c7100a318e46526da0d4c6d733b21df1d3dce1fc63c2108d84f19b079a00a06082a8648ce3d030107a144034200041104a6c56de39dfd25c8623e2b0e4d3f79866d756dddb4ecd78f09109a88333e3a26776136a0c59568ad311d44040ef876663f9d9d0e76e8cb71a40e9380e62d"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestVerify(t *testing.T) {
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	b, _ := x509.MarshalECPrivateKey(privKey)
	t.Logf("%x", b)
}

func TestSign(t *testing.T) {
	s := Sign("", *makeTestWallet())
}
