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
	TxID  string
	Index int
	Owner string
}

type TxOut struct {
	Owner  string
	Amount int
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("Not Enough Balance")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	// from의 txOut을 하나씩 뒤져서 그 tx의 amount의 총합이 amount와 클 때까지 합하고 그걸 실행하는 txIn의 amount를 쓴다. 잔돈처리는 후술.
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	// txIns -> {to, total}. txOuts -> {from, change}, {to, amount}
	change := total - amount
	if change != 0 {
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
