package db

import (
	"github.com/trusch/jamesd/packet"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) AddPacket(p *packet.Packet) error {
	c := db.session.DB("jamesd").C("packets")
	_, err := c.Upsert(bson.M{"name": p.Name, "tags": p.Tags}, p)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPacket(name string, tags []string) (*packet.Packet, error) {
	p := &packet.Packet{}
	c := db.session.DB("jamesd").C("packets")
	err := c.Find(bson.M{"name": name, "tags": bson.M{"$all": tags}}).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DB) RemovePacket(name string, tags []string) error {
	c := db.session.DB("jamesd").C("packets")
	err := c.Remove(bson.M{"name": name, "tags": tags})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ListPackets(name string, tags []string) ([]*packet.Packet, error) {
	c := db.session.DB("jamesd").C("packets")
	query := bson.M{}
	if name != "" {
		query["name"] = name
	}
	if len(tags) > 0 {
		query["tags"] = bson.M{"$all": tags}
	}
	packets := make([]*packet.Packet, 0)
	err := c.Find(query).Select(bson.M{"name": 1, "tags": 1}).Sort("name", "tags").All(&packets)
	if err != nil {
		return nil, err
	}
	return packets, nil
}
