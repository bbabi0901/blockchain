package p2p

import (
	"fmt"
	"net/http"

	"github.com/bbabi0901/blockchain/blockchain"
	"github.com/bbabi0901/blockchain/utils"
	"github.com/gorilla/websocket"
)

var conns []*websocket.Conn
var upgrader = websocket.Upgrader{}

// upgrade http connection(stateless) to WebSocket connection(stateful)
func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// 원래 CheckOrigin은 WebSocket 연결을 인증하여 true or false 반환 like cookie or session of HTTP. 일단은 항상 true 반환하도록 덮어쓰기.
	// Port :3000 will upgrade the request from :4000

	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	// RemoteAddr; 요청을 보낸 address를 기록 & 제공
	openPort := r.URL.Query().Get("openPort")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ip != "" && openPort != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	conns = append(conns, conn)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
}

func AddPeer(address, port, openPort string) {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port)
	sendNewestBlock(p)
}

func BroadcastNewBlock(b *blockchain.Block) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, peer := range Peers.v {
		notifyNewBlock(b, peer)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	Peers.m.Lock()
	defer Peers.m.Unlock()
	for _, peer := range Peers.v {
		notifyNewTx(tx, peer)
	}
}
