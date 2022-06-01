package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "30770201010420df9ea67c7100a318e46526da0d4c6d733b21df1d3dce1fc63c2108d84f19b079a00a06082a8648ce3d030107a144034200041104a6c56de39dfd25c8623e2b0e4d3f79866d756dddb4ecd78f09109a88333e3a26776136a0c59568ad311d44040ef876663f9d9d0e76e8cb71a40e9380e62d"
	testPayload string = "fd5853ee574381101ca6bd6962841ad10a6a9ffaf75e108027d4e9b394028dd4"
	testSig     string = "ac71515ac42ef1333c8d037dc932692f69231b0fc9df40d010adb63277feb69ef24ce3bbf21d6f02ff55f12d739400555a8c9c3a09723b4696de74b57322c5b6"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestWallet(t *testing.T) {
	t.Run("New Wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return false },
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New Wallet should return new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return true },
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("Wallet should return existing wallet instance")
		}
	})
}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestVerify(t *testing.T) {
	type testStruct struct {
		input string
		ok    bool
	}
	tests := []testStruct{
		{testPayload, true},
		{"fd5853ee574381101ca6bd6962841ad10a6a9ffaf75e108027d4e9b394028dd5", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, *makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s", s)
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex.")
	}
}
