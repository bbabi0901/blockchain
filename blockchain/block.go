package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/bbabi0901/blockchain/db"
	"github.com/bbabi0901/blockchain/utils"
)

type Block struct {
	Data      string `json:"data"`
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash,omitempty"`
	Height    int    `json:"height"`
	Diffculty int    `json:"difficulty"`
	Nonce     int    `json:"nonce"`
	Timestamp int    `json:"timestamp"`
}

var ErrNotFound = errors.New("Block not found")

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", Blockchain().difficulty())
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

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:      data,
		Hash:      "",
		PrevHash:  prevHash,
		Height:    height,
		Diffculty: Blockchain().difficulty(),
		Nonce:     0,
	}
	block.mine()
	block.persist()
	return block
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
