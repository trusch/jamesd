package server

import "github.com/trusch/jamesd/systemstate"

type MessageType int

const (
	// client side
	PING MessageType = iota
	STATE
	//server side
	PONG
	INSTALL
	UNINSTALL
	GET_STATE
	VPN_START
	VPN_STOP
)

type Message struct {
	Type   MessageType
	State  *systemstate.SystemState
	Packet []byte
}
