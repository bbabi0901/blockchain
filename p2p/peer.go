package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Peers map[string]*peer = make(map[string]*peer)

type peer struct {
	conn  *websocket.Conn
	inbox chan []byte
}

func (p *peer) read() {
	// delete peer in case of err
	for {
		_, m, err := p.conn.ReadMessage() // blocking for loop till it gets the message
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	// read inbox
	for {
		m := <-p.inbox // blocking for loop till inbox of peer gets the message
		err := p.conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			break
		}
	}
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	p := &peer{
		conn,
		make(chan []byte),
	}
	key := fmt.Sprintf("%s:%s", address, port)
	go p.read()
	go p.write()
	Peers[key] = p
	return p
}
