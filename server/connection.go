package server

import (
	"crypto/tls"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"sync"
)

type Connection struct {
	ID string

	conn    net.Conn
	mutex   sync.Mutex
	decoder *gob.Decoder
	encoder *gob.Encoder
	server  *Server
}

func NewConnection(conn net.Conn, server *Server) (*Connection, error) {
	connection := &Connection{
		conn:    conn,
		decoder: gob.NewDecoder(conn),
		encoder: gob.NewEncoder(conn),
		server:  server,
	}
	id, err := connection.getId()
	if err != nil {
		conn.Close()
		return nil, err
	}
	connection.ID = id
	go connection.readMessages()
	log.Printf("got new connection: %v", id)
	return connection, nil
}

func (conn *Connection) Send(msg *Message) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.encoder.Encode(msg)
}

func (conn *Connection) Read() (*Message, error) {
	msg := &Message{}
	err := conn.decoder.Decode(msg)
	return msg, err
}

func (conn *Connection) getId() (string, error) {
	if tlsConn, ok := conn.conn.(*tls.Conn); ok {
		err := tlsConn.Handshake()
		if err != nil {
			return "", err
		}
		cert := tlsConn.ConnectionState().PeerCertificates[0]
		return cert.Subject.CommonName, nil
	}
	return "", errors.New("not a TLS connection")
}

func (conn *Connection) readMessages() {
	for {
		msg, err := conn.Read()
		if err != nil {
			conn.handleDisconnect(err)
			break
		}
		conn.handleIncomingMessage(msg)
	}
}

func (conn *Connection) handleIncomingMessage(msg *Message) {
	switch msg.Type {
	case PING:
		{
			reply := &Message{Type: PONG}
			err := conn.Send(reply)
			if err != nil {
				conn.handleDisconnect(err)
				break
			}
			question := &Message{Type: GET_STATE}
			err = conn.Send(question)
			if err != nil {
				conn.handleDisconnect(err)
				break
			}
		}
	case STATE:
		{
			state := msg.State
			state.ID = conn.ID
			conn.server.handleNewState(state)
		}
	}
}

func (conn *Connection) handleDisconnect(err error) {
	log.Printf("client %v disconnected (%v)", conn.ID, err)
	conn.server.handleDisconnect(conn.ID)
	conn.conn.Close()
}
