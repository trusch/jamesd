package jamesd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"net"

	"gopkg.in/yaml.v2"
)

type Client struct {
	conn        net.Conn
	installRoot string
}

func (client *Client) GetPackage(request *Request) (*Response, error) {
	encoder := gob.NewEncoder(client.conn)
	decoder := gob.NewDecoder(client.conn)
	err := encoder.Encode(request)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	err = decoder.Decode(resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp, nil
}

func (client *Client) InstallPackage(response *Response) error {
	return InstallPackage(response, client.installRoot)
}

func NewClient(addr, certFile, keyFile, caFile, installRoot string) (*Client, error) {
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
	return &Client{conn, installRoot}, nil
}

type ClientConfig struct {
	Endpoint     string
	Cert         string
	Key          string
	CA           string
	InstallRoot  string
	Architecture string
	Services     []*ServiceConfig
}

type ServiceConfig struct {
	Name    string
	Version string
}

func (cfg *ClientConfig) ParseFile(file string) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, cfg)
}
