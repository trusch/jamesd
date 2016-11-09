package jamesd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
	db       *DB
}

func (server *Server) getIdFromConn(conn net.Conn) (string, error) {
	if tlsConn, ok := conn.(*tls.Conn); ok {
		err := tlsConn.Handshake()
		if err != nil {
			return "", err
		}
		cert := tlsConn.ConnectionState().PeerCertificates[0]
		return cert.Subject.CommonName, nil
	}
	return "", errors.New("not a TLS connection")
}

func (server *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	clientId, err := server.getIdFromConn(conn)
	if err != nil {
		log.Print(err)
		return
	}
	defer func() { log.Print("lost connection to ", clientId) }()
	log.Print("new connection from ", clientId)
	for {
		decoder := gob.NewDecoder(conn)
		request := &Request{}
		err = decoder.Decode(request)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			break
		}
		log.Printf("%v requests %v_%v:%v", clientId, request.Package, request.Arch, request.Version)
		var response *Response
		if !server.db.CheckPermission(clientId, request.Package) {
			log.Printf("%v's request for %v_%v:%v could not be fullfilled: unauthorized", clientId, request.Package, request.Arch, request.Version)
			response = &Response{Error: "unauthorized"}
		} else {
			response = GetPackage(request, server.db)
		}
		if response.Error != "" {
			log.Printf("%v's request for %v_%v:%v could not be fullfilled: %v", clientId, request.Package, request.Arch, request.Version, response.Error)
		} else {
			log.Printf("send %v_%v:%v to %v", request.Package, request.Arch, request.Version, clientId)
		}
		encoder := gob.NewEncoder(conn)
		err = encoder.Encode(response)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			break
		}
	}
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

func NewServer(addr, certFile, keyFile, caFile, dbUrl string) *Server {
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
	db, err := NewDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{ln, db}
	return server
}
