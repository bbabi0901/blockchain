package blockchain

import (
	"errors"
	"sync"
	"time"

	"github.com/bbabi0901/blockchain/utils"
	"github.com/bbabi0901/blockchain/wallet"
)

const minerReward int = 50

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	TxIns     []*TxIn  `json:"txInputs"`
	TxOuts    []*TxOut `json:"txOutputs"`
}

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
	Txs map[string]*Tx
	M   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{
			Txs: make(map[string]*Tx),
		}
	})
	return m
}

var ErrorNoMoney = errors.New("Not enough funds")
var ErrorNotValid = errors.New("Transaction Invalid")

func (m *mempool) TxToConfirm() []*Tx {
	var txs []*Tx
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs = append(txs, coinbase)
	for _, tx := range m.Txs {
		txs = append(txs, tx)
	}
	m.Txs = make(map[string]*Tx) // empty mempool
	return txs
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.M.Lock()
	defer m.M.Unlock()
	m.Txs[tx.ID] = tx
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return nil, err
	}
	m.Txs[tx.ID] = tx
	return tx, nil
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, *wallet.Wallet())
	}
}

// to check if unspent tx output is on mempool.
func isOnMempool(uTxOut *UTxOut) bool {
	exists := false

Outer: // labeling each for loops
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			exists = input.TxID == uTxOut.TxID && input.Index == uTxOut.Index
			break Outer // you can only break inner loop if there's no labeling
		}
	}
	return exists
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
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
		ID:        "",
		Timestamp: time.Now().Unix(),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

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
		ID:        "",
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
