package p2p

import (
	"fmt"
	"net/http"
	"time"

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
	peer := initPeer(conn, ip, openPort)

	time.Sleep(20 * time.Second)
	peer.inbox <- []byte("Hello from Port 3000!") // unblock for loop inside the wrtie() function of peer
}

func AddPeer(address, port, openPort string) {
	// Port :4000 is requesting an upgrade from the port :3000
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:]), nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)

	time.Sleep(10 * time.Second)
	peer.inbox <- []byte("Hello from Port 4000!") // unblock for loop inside the wrtie() function of peer
}
