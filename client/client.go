package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"github.com/trusch/jamesd/installer"
	"github.com/trusch/jamesd/server"
	"github.com/trusch/jamesd/systemstate"
)

type Client struct {
	conn        net.Conn
	installRoot string
	mutex       sync.Mutex
	encoder     *gob.Encoder
	decoder     *gob.Decoder
	state       *systemstate.SystemState
	stateFile   string
}

func (cli *Client) Send(msg *server.Message) error {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()
	return cli.encoder.Encode(msg)
}

func (cli *Client) Read() (*server.Message, error) {
	msg := &server.Message{}
	err := cli.decoder.Decode(msg)
	return msg, err
}

func New(addr, certFile, keyFile, caFile, installRoot, systemStateFile string) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)
	config := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ServerName:   "jamesd",
	}

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	state := &systemstate.SystemState{}
	err = state.Load(systemStateFile)
	if err != nil {
		log.Print("failed loading systemstate from disk")
	}
	state.Save(systemStateFile)

	cli := &Client{
		conn:        conn,
		installRoot: installRoot,
		encoder:     gob.NewEncoder(conn),
		decoder:     gob.NewDecoder(conn),
		stateFile:   systemStateFile,
		state:       state,
	}

	go cli.sendPings()

	return cli, nil
}

func (cli *Client) Run() {
	for {
		msg, err := cli.Read()
		if err != nil {
			log.Print("Error: ", err)
			break
		}
		cli.handleIncomingMessage(msg)
	}
}

func (cli *Client) handleIncomingMessage(msg *server.Message) {
	switch msg.Type {
	case server.PONG:
		{
			log.Print("got pong from server.")
		}
	case server.INSTALL:
		{
			pack := msg.Packet
			tar, err := pack.GetTarReader()
			if err != nil {
				log.Print("Error: ", err)
				break
			}
			err = installer.Install(tar, cli.installRoot, pack.PreInstallScript, pack.PostInstallScript)
			if err != nil {
				log.Print("Error: ", err)
				break
			}
			cli.state.MarkAppInstalled(&systemstate.AppInfo{Name: pack.Name, Tags: pack.Tags})
			log.Printf("installed %v:%v", pack.Name, pack.Tags)
			err = cli.state.Save(cli.stateFile)
			if err != nil {
				log.Print("failed writing statefile: ", err)
			}
			reply := &server.Message{
				Type:  server.STATE,
				State: cli.state,
			}
			err = cli.Send(reply)
			if err != nil {
				log.Print("Error: ", err)
				break
			}
		}
	case server.UNINSTALL:
		{
			pack := msg.Packet
			tar, err := pack.GetTarReader()
			if err != nil {
				log.Print("Error: ", err)
				break
			}
			err = installer.Uninstall(tar, cli.installRoot, pack.PreInstallScript, pack.PostInstallScript)
			if err != nil {
				log.Print("Error: ", err)
				break
			}
			cli.state.MarkAppUninstalled(&systemstate.AppInfo{Name: pack.Name, Tags: pack.Tags})
			log.Printf("uninstalled %v:%v", pack.Name, pack.Tags)
			err = cli.state.Save(cli.stateFile)
			if err != nil {
				log.Print("failed writing statefile: ", err)
			}
			msg := &server.Message{
				Type:  server.STATE,
				State: cli.state,
			}
			err = cli.Send(msg)
			if err != nil {
				log.Print("Error: ", err)
				break
			}
		}
	case server.GET_STATE:
		{
			msg := &server.Message{
				Type:  server.STATE,
				State: cli.state,
			}
			err := cli.Send(msg)
			if err != nil {
				log.Print("Error: ", err)
				break
			}
		}
	}
}

func (cli *Client) sendPings() {
	msg := &server.Message{Type: server.PING}
	for {
		time.Sleep(30 * time.Second)
		err := cli.Send(msg)
		if err != nil {
			log.Print("failed sending ping")
		}
	}
}
