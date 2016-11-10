package db

import (
	"github.com/trusch/jamesd/packet"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) AddPacket(p *packet.Packet) error {
	c := db.session.DB("jamesd").C("packets")
	_, err := c.Upsert(bson.M{"name": p.Name, "arch": p.Arch, "version": p.Version}, p)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPacket(name, arch, version string) (*packet.Packet, error) {
	p := &packet.Packet{}
	c := db.session.DB("jamesd").C("packets")
	err := c.Find(bson.M{"name": name, "arch": arch, "version": version}).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DB) RemovePacket(name, arch, version string) error {
	c := db.session.DB("jamesd").C("packets")
	err := c.Remove(bson.M{"name": name, "arch": arch, "version": version})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ListPackets() ([]*packet.Packet, error) {
	c := db.session.DB("jamesd").C("packets")
	packets := make([]*packet.Packet, 0)
	err := c.Find(nil).Select(bson.M{"name": 1, "version": 1, "arch": 1}).All(&packets)
	if err != nil {
		return nil, err
	}
	return packets, nil
}
