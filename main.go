package main

import (
	"github.com/bbabi0901/blockchain/cli"
	"github.com/bbabi0901/blockchain/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
