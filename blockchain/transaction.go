package blockchain

import (
	"errors"
	"time"

	"github.com/bbabi0901/blockchain/utils"
	"github.com/bbabi0901/blockchain/wallet"
)

const minerReward int = 50

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txInputs"`
	TxOuts    []*TxOut `json:"txOutputs"`
}

// get TxIn from UTxOut
type TxIn struct {
	TxID      string `json:"txId"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Amount int
	Index  int
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

// to check if unspent tx output is on mempool.
func isOnMempool(uTxOut *UTxOut) bool {
	exists := false

Outer: // labeling the loop
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			exists = input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
			break Outer // you can only break inner loop if there's no labeling
		}
	}
	return exists
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.Id, *wallet.Wallet())
	}
}

// verifying wheter we own the tx output, which is being used for tx input
// output of prevTx로부터 address를 가져오고 이를 새로운 tx input의 서명을 검증하는데 이용.
func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.Id, address) // payload, what is to verify, is transaction ID
		if !valid {
			break
		}
	}
	return valid
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: time.Now().Unix(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

// add transaction to mempool before adding it to block
func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

var ErrorNoMoney = errors.New("Not enough funds")
var ErrorNotValid = errors.New("Tx Invalid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())

	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: time.Now().Unix(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

// empty the memory pool, return the transaction
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
