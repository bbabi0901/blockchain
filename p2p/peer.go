package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// To protect map from data race, make Peers as struct instead of variable.
// And then put map inside the struct to protect, create mutex <- block the struct until unlocking
type peers struct {
	v map[string]*peer
	m sync.Mutex // Mutex will lock the struct where it is
}

var Peers peers = peers{
	v: make(map[string]*peer),
}

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

func (p *peer) close() {
	Peers.m.Lock()         // Lock peers so that other go routines can't modify before unlocking.
	defer Peers.m.Unlock() // Unlock the peers after closing
	p.conn.Close()
	delete(Peers.v, p.key)
}

func (p *peer) read() {
	// defer is the code that runs after the function has finished -> if loop break, the function ends. Then, go runs p.close()
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m) // blocking for loop till it gets the message
		if err != nil {
			break
		}
		handleMessage(&m, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox // blocking for loop till inbox of peer gets the message
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func AllPeers(p *peers) []string {
	p.m.Lock()
	defer p.m.Unlock()
	var keys []string

	for key := range p.v {
		keys = append(keys, key)
	}
	return keys
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.m.Lock()
	defer Peers.m.Unlock()

	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		address: address,
		port:    port,
	}
	go p.read()
	go p.write()
	Peers.v[key] = p
	return p
}
