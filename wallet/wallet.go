package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/bbabi0901/blockchain/utils"
)

const (
	fileName string = "BBaBi.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644) // 0644 -> permission for read and write
	utils.HandleErr(err)
}

// named return
func restoreKey() (privKey *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	privKey, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

// making address from x and y in publicKey
func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func Sign(payload string, w wallet) string {
	// string to bytes
	pAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, pAsBytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func decodeBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	// 받아온 signature로부터 r, s 만들어서 payload로부터 생성된 것과 비교, 검증
	r, s, err := decodeBigInts(signature)
	utils.HandleErr(err)

	// 받아온 address로부터 publickey를 만들어서 비교, 검증
	x, y, err := decodeBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	pAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&publicKey, pAsBytes, r, s)

	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
