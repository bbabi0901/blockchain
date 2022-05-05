package blockchain

import (
	"fmt"
	"sync"

	"github.com/bbabi0901/blockchain/db"
	"github.com/bbabi0901/blockchain/utils"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			// search for checkpoint on the db, then restore the blockchain from the bytes
			// so gonna make func in db.go to find data from db
			checkpoint := db.Checkpoint()
			if checkpoint == nil { // db의 checkpoint에 저장된 데이터가 없으면 만들어 놓은 empty blockchain에 최초의 block생성
				b.AddBlock("Genesis Block")
			} else { // db의 checkpoint에 저장된 데이터가 있으면 byte로 저장된 blockchain을 decoding해서 복원
				b.restore(checkpoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
