package main

import (
	"github.com/bbabi0901/blockchain/cli"
	"github.com/bbabi0901/blockchain/db"
)

func main() {
	defer db.Close() // defer -> 함수가 종료될 때 실행된다 -> main function 종료시 database도 close...

	cli.Start()
}
