package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/spec"
	"github.com/trusch/jamesd/state"
)

// Client is a jamesd client class
type Client struct {
	endpoint string
	client   *http.Client
}

// New returns a new client
func New(endpoint string) *Client {
	return &Client{endpoint, &http.Client{}}
}

// GetPackets returns a list of all packet control infos
func (cli *Client) GetPackets() (map[string][]*packet.ControlInfo, error) {
	req, err := http.NewRequest("GET", cli.endpoint+"/packet/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	result := make(map[string][]*packet.ControlInfo)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UploadPacket sends a packet to the server
func (cli *Client) UploadPacket(pack *packet.Packet) error {
	data, err := pack.ToData()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", cli.endpoint+"/packet/", bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	return nil
}

// GetDesiredState returns a list of packets to be installed for a given labelset
func (cli *Client) GetDesiredState(labels map[string]string) (*state.State, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.Encode(labels)
	req, err := http.NewRequest("POST", cli.endpoint+"/packet/compute", buf)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	result := &state.State{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeletePacket deletes a packet from the server
func (cli *Client) DeletePacket(hash string) error {
	url := fmt.Sprintf("%v/packet/%v", cli.endpoint, hash)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	return nil
}

// GetPacketInfo returns the packet info to a given hash
func (cli *Client) GetPacketInfo(hash string) (*packet.ControlInfo, error) {
	url := fmt.Sprintf("%v/packet/%v/info", cli.endpoint, hash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	decoder := json.NewDecoder(resp.Body)
	info := &packet.ControlInfo{}
	err = decoder.Decode(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetPacketData returns the packet info to a given hash
func (cli *Client) GetPacketData(hash string) (*packet.Packet, error) {
	url := fmt.Sprintf("%v/packet/%v/data", cli.endpoint, hash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pack, err := packet.NewFromData(bs)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

// GetSpecs returns a list of all packet control infos
func (cli *Client) GetSpecs() ([]*spec.Spec, error) {
	req, err := http.NewRequest("GET", cli.endpoint+"/spec/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	result := []*spec.Spec{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UploadSpec sends a packet to the server
func (cli *Client) UploadSpec(s *spec.Spec) error {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", cli.endpoint+"/spec/", buf)
	if err != nil {
		return err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	return nil
}

// GetMergedSpec returns the merged specs for a given labelset
func (cli *Client) GetMergedSpec(labels map[string]string) (*spec.Spec, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(labels)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", cli.endpoint+"/spec/compute", buf)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	result := &spec.Spec{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetSpec returns a specific spec
func (cli *Client) GetSpec(id string) (*spec.Spec, error) {
	req, err := http.NewRequest("GET", cli.endpoint+"/spec/"+id, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	result := &spec.Spec{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PutSpec sends a packet to the server
func (cli *Client) PutSpec(s *spec.Spec) error {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", cli.endpoint+"/spec/"+s.ID, buf)
	if err != nil {
		return err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	return nil
}

// DeleteSpec returns a specific spec
func (cli *Client) DeleteSpec(id string) error {
	req, err := http.NewRequest("DELETE", cli.endpoint+"/spec/"+id, nil)
	if err != nil {
		return err
	}
	resp, err := cli.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return errors.New("http error: " + strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}
	return nil
}
