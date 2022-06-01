package db

import (
	"fmt"
	"os"

	"github.com/bbabi0901/blockchain/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName       = "blockchain"
	dataBucket   = "data" // Bucket is like table of SQL
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

type DB struct{}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}

func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}

func (DB) SaveChain(data []byte) {
	saveChain(data)
}

func (DB) LoadChain() []byte {
	return loadChain()
}

func (DB) DeleteAllBlocks() {
	emptyBlocks()
}

func getDbName() string {
	port := os.Args[2][6:]
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

func InitDB() {
	if db == nil {
		// init db
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
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
}

func Close() {
	db.Close()
}

func saveBlock(hash string, data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func saveChain(data []byte) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func loadChain() []byte {
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint)) // Get()은 checkpoint라는 key에 값이 있으면 해당 값을 없으면 nil을 반환
		return nil
	})
	return data
}

func findBlock(hash string) []byte {
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}

func emptyBlocks() {
	db.Update(func(tx *bolt.Tx) error {
		utils.HandleErr(tx.DeleteBucket([]byte(blocksBucket)))
		_, err := tx.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
