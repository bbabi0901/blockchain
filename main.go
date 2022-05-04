package main

import (
	"fmt"

	"github.com/bbabi0901/blockchain/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Second Block")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("%s\n", block.Data)
		fmt.Printf("%s\n", block.Hash)
		fmt.Printf("%s\n", block.PrevHash)
	}

	fmt.Println("-----------")

	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")
	for _, block := range chain.AllBlocks() {
		fmt.Printf("%s\n", block.Data)
		fmt.Printf("%s\n", block.Hash)
		fmt.Printf("%s\n", block.PrevHash)
	}

}
