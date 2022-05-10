package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/bbabi0901/blockchain/utils"
)

const (
	hashedMessage string = "7509dea1e30fc3c3346abee17b2559b22948e46a7074db7fa45080b85174c039"
	privateKey    string = "307702010104206f88460dde0830a85d338201e8078a4b34c71c670de17889f622e38288f41daaa00a06082a8648ce3d030107a14403420004e2a0884dc9d5897741c3c1519ce3ae9ca97e8c6c2139da336d4f70a9991f0ae3f1de15bb5d676c352c52bc37963cef1bedeeba746211fb974256234e1d424dfd"
	signature     string = "0aa87d97f078c43f430952b8f1d189bb4eef1c7d71c09856d7efedcb220f8b081d8bcbf10867cb649f9bb7adbc3c12f942bea7dc51dccd103a7bed37167288c2"
)

func Start() {
	// just to be sure if the encoding of the privateKey is hexadecimal and the format is correct
	privateBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	private, err := x509.ParseECPrivateKey(privateBytes)
	utils.HandleErr(err)

	signatureBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)

	rBytes := signatureBytes[:len(signatureBytes)/2]
	sBytes := signatureBytes[len(signatureBytes)/2:]

	// initaillizing *big.Int
	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	hashBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)
	fmt.Println(ok)
}
