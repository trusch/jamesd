package server

import (
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/systemstate"
)

type MessageType int

const (
	// client side
	PING MessageType = iota
	STATE
	//server side
	PONG
	INSTALL
	GET_STATE
	VPN_START
	VPN_STOP
)

type Message struct {
	Type   MessageType
	State  *systemstate.SystemState
	Packet *packet.Packet
}
