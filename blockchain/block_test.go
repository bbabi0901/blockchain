package blockchain

import (
	"reflect"
	"testing"

	"github.com/bbabi0901/blockchain/utils"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block Not Found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				return nil
			},
		}
		_, err := FindBlock("x")
		if err == nil {
			t.Error("The block should not be found.")
		}
	})
	t.Run("Block is Found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				b := &Block{
					Height: 1,
				}
				return utils.ToBytes(b)
			},
		}
		block, _ := FindBlock("x")
		if block.Height != 1 {
			t.Error("The block should be found.")
		}
	})
}
