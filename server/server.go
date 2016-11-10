package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"

	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/systemstate"
)

type Server struct {
	listener    net.Listener
	db          *db.DB
	connections map[string]*Connection
}

func (server *Server) handleDisconnect(id string) {
	delete(server.connections, id)
}

func (server *Server) handleNewState(currentState *systemstate.SystemState) {
	err := server.db.SaveCurrentSystemState(currentState)
	if err != nil {
		log.Print("Error: ", err)
	}
	desiredState, err := server.db.GetDesiredSystemState(currentState.ID)
	if err != nil {
		log.Print("Error: ", err)
		return
	}
	neededApps := make([]*systemstate.AppInfo, 0)
	for _, desiredApp := range desiredState.Apps {
		isNeeded := true
		for _, currentApp := range currentState.Apps {
			if currentApp.Name == desiredApp.Name && currentApp.Version == desiredApp.Version {
				isNeeded = false
				break
			}
		}
		if isNeeded {
			neededApps = append(neededApps, desiredApp)
		}
	}
	for _, appInfo := range neededApps {
		pack, e := server.db.GetPacket(appInfo.Name, currentState.Arch, appInfo.Version)
		if e != nil {
			log.Print("Error: ", e)
			return
		}
		msg := &Message{
			Type:   INSTALL,
			Packet: pack,
		}
		e = server.connections[currentState.ID].Send(msg)
		if e != nil {
			log.Print("Error: ", e)
			return
		}
	}
}

func (server *Server) handleConn(conn net.Conn) {
	connection, err := NewConnection(conn, server)
	if err != nil {
		log.Print(err)
		return
	}
	msg := &Message{Type: GET_STATE}
	connection.Send(msg)
	server.connections[connection.ID] = connection
}

func (server *Server) Run() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go server.handleConn(conn)
	}
}

func New(addr, certFile, keyFile, caFile, dbUrl string) *Server {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)
	config := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
	}
	ln, err := tls.Listen("tcp", addr, config)
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.New(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{ln, database, make(map[string]*Connection)}
	return server
}
