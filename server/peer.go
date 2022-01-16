package main

import (
	"github.com/atrox/haikunatorgo/v2"
	"github.com/gorilla/websocket"
	"github.com/mssola/user_agent"
	"hash/fnv"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Peer struct {
	Conn         *websocket.Conn
	IP           string
	ID           string
	RtcSupported bool
	Timer        *time.Timer
	LastBeat     time.Time
	Name         PeerName
}

type PeerName struct {
	Model       string `json:"model"`
	OS          string `json:"os"`
	Browser     string `json:"browser"`
	Type        string `json:"type"`
	DeviceName  string `json:"deviceName"`
	DisplayName string `json:"displayName"`
}

type PeerInfo struct {
	ID           string   `json:"id"`
	Name         PeerName `json:"name"`
	RtcSupported bool     `json:"rtcSupported"`
}

func newPeer(conn *websocket.Conn, req *http.Request) *Peer {

	peer := Peer{Conn: conn}

	peer.SetIP(req)
	peer.SetID(req)

	peer.RtcSupported = strings.Contains(req.URL.Path, "webrtc")
	peer.setName(req)

	peer.LastBeat = time.Now()

	return &peer
}

func (peer *Peer) SetIP(req *http.Request) {

	ip, err := getIP(req)
	if err != nil {
		log.Println(err)
	}

	peer.IP = ip

	// IPv4 and IPv6 use different values to refer to localhost
	if peer.IP == "::1" || peer.IP == "::ffff:127.0.0.1" {
		peer.IP = "127.0.0.1"
	}
}

func (peer *Peer) SetID(req *http.Request) {

	peerIDCookie, err := req.Cookie("peerid")
	if err == nil {
		peer.ID = peerIDCookie.Value
	}
}

func (peer *Peer) GetIDHash() int64 {
	h := fnv.New64a()
	h.Write([]byte(peer.ID))
	return int64(h.Sum64())
}

func (peer *Peer) getInfo() PeerInfo {
	return PeerInfo{ID: peer.ID, Name: peer.Name, RtcSupported: peer.RtcSupported}
}

func (peer *Peer) setName(req *http.Request) {
	var ua = user_agent.New(req.Header.Get("USER-AGENT"))

	deviceName := ""

	var OS = ua.OS()

	if OS != "" {
		deviceName = strings.Replace(OS, "Mac OS", "Mac", -1) + " "
	}

	browserName, _ := ua.Browser()
	deviceName += browserName

	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	var dtype = ""
	if ua.Mobile() {
		dtype = "mobile"
	} else {
		dtype = "other"
	}

	haikunator := haikunator.New()
	haikunator.TokenLength = 0
	haikunator.Delimiter = " "
	haikunator.Random = rand.New(rand.NewSource(peer.GetIDHash()))

	peer.Name = PeerName{
		Model:       browserName,
		OS:          OS,
		Browser:     browserName,
		Type:        dtype,
		DeviceName:  deviceName,
		DisplayName: strings.Title(haikunator.Haikunate()),
	}
}
