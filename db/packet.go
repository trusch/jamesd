package db

import (
	"github.com/trusch/jamesd/packet"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) AddPacket(p *packet.Packet) error {
	c := db.session.DB("jamesd").C("packets")
	_, err := c.Upsert(bson.M{"controlinfo.name": p.Name, "controlinfo.tags": p.Tags}, p)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPacket(name string, tags []string) (*packet.Packet, error) {
	p := &packet.Packet{}
	c := db.session.DB("jamesd").C("packets")
	err := c.Find(bson.M{"controlinfo.name": name, "controlinfo.tags": tags}).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DB) RemovePacket(name string, tags []string) error {
	c := db.session.DB("jamesd").C("packets")
	err := c.Remove(bson.M{"controlinfo.name": name, "controlinfo.tags": tags})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ListPackets(name string, tags []string) (packet.PacketList, error) {
	c := db.session.DB("jamesd").C("packets")
	query := bson.M{}
	if name != "" {
		query["controlinfo.name"] = name
	}
	if len(tags) > 0 {
		query["controlinfo.tags"] = bson.M{"$all": tags}
	}
	packets := packet.PacketList{}
	err := c.Find(query).Select(bson.M{"controlinfo.name": 1, "controlinfo.tags": 1}).Sort("controlinfo.name", "controlinfo.tags").All(&packets)
	if err != nil {
		return nil, err
	}
	return packets, nil
}

func (db *DB) GetMatchingPackets(name string, tags []string) (packet.PacketList, error) {
	c := db.session.DB("jamesd").C("packets")
	// -> give all docs where in controlinfo.tags is NOT an element which is NOT in the query-tag-set
	query := bson.M{
		"controlinfo.tags": bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					"$nin": tags,
				},
			},
		},
	}
	if name != "" {
		query["controlinfo.name"] = name
	}
	packets := packet.PacketList{}
	err := c.Find(query).All(&packets)
	if err != nil {
		return nil, err
	}
	return packets, nil
}
