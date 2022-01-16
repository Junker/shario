package main

import (
	"flag"
	"github.com/Jeffail/gabs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	addr     = flag.String("addr", "localhost:3000", "server address")
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
	}

	rooms = make(map[string]map[string]*Peer)

	pingMessage = newMessagePing()
)

func onConnect(w http.ResponseWriter, req *http.Request) {

	var cookieHeader = http.Header{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	_, err := req.Cookie("peerid")

	if err != nil {
		uuid, _ := generateUUID()

		cookieHeader.Set("Set-Cookie", "peerid="+uuid+"; SameSite=Strict; Secure")
	}

	conn, err := upgrader.Upgrade(w, req, cookieHeader)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer conn.Close()

	peer := newPeer(conn, req)
	// fmt.Printf("PEER: %+v", peer)

	joinRoom(peer)
	keepAlive(peer)

	// send displayName
	sendJSON(peer, newMessageDisplayName(MessageData{DisplayName: peer.Name.DisplayName, DeviceName: peer.Name.DeviceName}))

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		onMessage(peer, message)
	}

}

func onMessage(sender *Peer, msgData []byte) {

	msgJSON, err := gabs.ParseJSON(msgData)
	if err != nil {
		log.Println(err)
		return
	}

	messageType, _ := msgJSON.Path("type").Data().(string)
	messageTo, hasMessageTo := msgJSON.Path("to").Data().(string)

	if messageType == "disconnect" {
		leaveRoom(sender)
	}
	if messageType == "pong" {
		sender.LastBeat = time.Now()
	}

	// relay message to recipient
	_, roomExists := rooms[sender.IP]
	if hasMessageTo && roomExists {
		var recipientID = messageTo
		var recipient = rooms[sender.IP][recipientID]

		msgJSON.SetP("", "to")
		msgJSON.SetP(sender.ID, "sender")

		sendBytes(recipient, msgJSON.Bytes())

		return
	}
}

func joinRoom(peer *Peer) {
	log.Printf("Peer %s joined the room with IP %s", peer.ID, peer.IP)

	var otherPeers, roomExists = rooms[peer.IP]

	// notify all other peers
	if roomExists {
		for _, otherPeer := range otherPeers {
			message := newMessagePeerJoined(peer.getInfo())
			sendJSON(otherPeer, message)
		}
	} else {
		// if room doesn't exist, create it
		rooms[peer.IP] = make(map[string]*Peer)
	}

	// notify peer about the other peers
	otherPeerInfos := []PeerInfo{}
	for _, otherPeer := range otherPeers {
		otherPeerInfos = append(otherPeerInfos, otherPeer.getInfo())
	}
	sendJSON(peer, newMessagePeers(otherPeerInfos))

	// add peer to room
	rooms[peer.IP][peer.ID] = peer
}

func leaveRoom(peer *Peer) {
	log.Printf("Peer %s leaved the room with IP %s", peer.ID, peer.IP)

	var _, peerIn = rooms[peer.IP][peer.ID]
	if !peerIn {
		return
	}

	cancelKeepAlive(rooms[peer.IP][peer.ID])

	// delete the peer
	delete(rooms[peer.IP], peer.ID)

	peer.Conn.Close()

	//if room is empty, delete the room
	if len(rooms[peer.IP]) == 0 {
		delete(rooms, peer.IP)
	} else {
		for _, otherPeer := range rooms[peer.IP] {
			sendJSON(otherPeer, newMessagePeerLeft(peer.ID))
		}
	}
}

func keepAlive(peer *Peer) {
	cancelKeepAlive(peer)
	var timeout = 15 * time.Second

	if time.Now().Sub(peer.LastBeat) > 4*timeout {
		leaveRoom(peer)
		return
	}

	sendJSON(peer, pingMessage)

	peer.Timer = time.AfterFunc(timeout, func() { keepAlive(peer) })
}

func cancelKeepAlive(peer *Peer) {
	if peer.Timer != nil {
		peer.Timer.Stop()
	}
}

func sendJSON(peer *Peer, msgdata interface{}) {
	peer.Conn.WriteJSON(msgdata)
}

func sendBytes(peer *Peer, msgdata []byte) {
	peer.Conn.WriteMessage(websocket.TextMessage, msgdata)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", onConnect)

	log.Printf("Server listening on addr %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
