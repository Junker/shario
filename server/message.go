package main

type MessageDisplayName struct {
	Type    string      `json:"type""`
	Message MessageData `json:"message"`
}

type MessagePeerLeft struct {
	Type   string `json:"type"`
	PeerID string `json:"peerId"`
}

type MessagePeers struct {
	Type  string     `json:"type"`
	Peers []PeerInfo `json:"peers"`
}

type MessagePeerJoined struct {
	Type string   `json:"type"`
	Peer PeerInfo `json:"peer"`
}

type MessagePing struct {
	Type string `json:"type"`
}

type MessageData struct {
	DisplayName string `json:"displayName"`
	DeviceName  string `json:"deviceName"`
}

func newMessageDisplayName(msgData MessageData) MessageDisplayName {
	return MessageDisplayName{Type: "display-name", Message: msgData}
}

func newMessagePeerLeft(peerID string) MessagePeerLeft {
	return MessagePeerLeft{Type: "peer-left", PeerID: peerID}
}

func newMessagePeers(peers []PeerInfo) MessagePeers {
	return MessagePeers{Type: "peers", Peers: peers}
}

func newMessagePeerJoined(peer PeerInfo) MessagePeerJoined {
	return MessagePeerJoined{Type: "peer-joined", Peer: peer}
}

func newMessagePing() MessagePing {
	return MessagePing{Type: "ping"}
}
