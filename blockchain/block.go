package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/bbabi0901/blockchain/db"
	"github.com/bbabi0901/blockchain/utils"
)

type Block struct {
	Hash        string `json:"hash"`
	PrevHash    string `json:"prevHash,omitempty"`
	Height      int    `json:"height"`
	Diffculty   int    `json:"difficulty"`
	Nonce       int    `json:"nonce"`
	Timestamp   int    `json:"timestamp"`
	Transaction []*Tx  `json:"transaction"`
}

var ErrNotFound = errors.New("Block not found")

func (b *Block) mine() {
	target := strings.Repeat("0", getDifficulty(Blockchain()))
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func persist(b *Block) {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func restore(data []byte, b *Block) {
	utils.FromBytes(b, data)
}

func createBlock(prevHash string, height, diff int) *Block {
	block := &Block{
		Hash:        "",
		PrevHash:    prevHash,
		Height:      height,
		Diffculty:   diff,
		Nonce:       0,
		Transaction: []*Tx{makeCoinbaseTx("BBaBi")},
	}
	block.mine()
	block.Transaction = Mempool.TxToConfirm()
	persist(block)
	return block
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	restore(blockBytes, block)
	return block, nil
}
