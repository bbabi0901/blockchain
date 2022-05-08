package blockchain

import (
	"sync"

	"github.com/bbabi0901/blockchain/db"
	"github.com/bbabi0901/blockchain/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		// search for checkpoint on the db, then restore the blockchain from the bytes
		// so gonna make func in db.go to find data from db
		checkpoint := db.Checkpoint()
		if checkpoint == nil { // db의 checkpoint에 저장된 데이터가 없으면 만들어 놓은 empty blockchain에 최초의 block생성
			b.AddBlock()
		} else { // db의 checkpoint에 저장된 데이터가 있으면 byte로 저장된 blockchain을 decoding해서 복원
			b.restore(checkpoint)
		}
	})
	return b
}

/*
go rule:
Method is only used when it’s mutating the struct
If not, we use function.
Modifying -> method
Just using struct as an input -> function
*/

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Diffculty
	persistBlockchain(b)
}

// mutating the difficulty of *blockchain, thus it should be method
func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

// mutating the difficulty of *blockchain, thus it should be method
func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		// recalculate difficulty
		return recalculateDiffculty(b)
	} else {
		return b.CurrentDifficulty
	}
}

// Doesn't change *blockchain but just take it as an input for newestHash. It should be a function
func recalculateDiffculty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp - lastRecalculatedBlock.Timestamp) / 60 // min
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

// Doesn't change *blockchain but just take it as an input for newestHash. It should be a function
func persistBlockchain(b *blockchain) {
	db.SaveCheckpoint(utils.ToBytes(b))
}

// Doesn't change *blockchain but just take it as an input for newestHash. It should be a function
func Blocks(b *blockchain) []*Block {
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

// Unspent transaction
func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	// finding tx output that hasn't been referenced
	for _, block := range Blocks(b) {
		for _, tx := range block.Transaction {
			// marking tx id to track the output that are being used to create input -> marked output = spent output
			for _, input := range tx.TxIns {
				if input.Owner == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Owner == address {
					// if boolean is not TRUE, it means the output is not marked. That means it is unspent output
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UTxOut{tx.Id, output.Amount, index}
						// if unspent tx out is already on mempool, there's no need to append tx again.
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}
