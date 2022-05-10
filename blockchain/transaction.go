package blockchain

import (
	"errors"
	"time"

	"github.com/bbabi0901/blockchain/utils"
)

const minerReward int = 50

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txInputs"`
	TxOuts    []*TxOut `json:"txOutputs"`
}

type TxIn struct {
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
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
	tx, err := makeTx("BBaBi", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, errors.New("Not enough funds")
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

	return tx, nil
}

// empty the memory pool, return the transaction
func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("BBaBi")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
