package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/bbabi0901/blockchain/blockchain"
	"github.com/bbabi0901/blockchain/utils"
)

type MessageKind int

const (
	MessageNewestBlock       MessageKind = iota // All variable under "iota" will get each number and same kind of iota
	MessageAllBlocksRequest                     // First variable is 0, and this will be 1 bco iota
	MessageAllBlocksResponse                    // this will be 2 bco iota
	MessageNewBlockNotify
	MessageNewTxNotify
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}
	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, b)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		// if the function is running by port 3000, b is the block of port 3000 and p is that of port 4000
		// so if block height of port 4000 is higher than that of port 3000, port 3000 will request all the blocks from 4000
		// else, port 3000 will send all the blocks to 4000
		if payload.Height >= b.Height {
			// b-side will request all the blocks from p-side
			requestAllBlocks(p)
		} else {
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s wants all the blocks.\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all blocks from %s.\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	}
}
