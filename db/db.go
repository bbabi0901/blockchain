package db

import (
	"github.com/bbabi0901/blockchain/utils"
	"github.com/boltdb/bolt"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data" // Bucket is like table of SQL
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		// init db
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		// read-write transaction by using db.Update()
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error { // tx: transaction
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket)) // returns bucket and error
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint)) // Get()은 checkpoint라는 key에 값이 있으면 해당 값을 없으면 nil을 반환
		return nil
	})
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}
