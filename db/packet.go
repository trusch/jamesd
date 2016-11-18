package db

import (
	"github.com/trusch/jamesd/packet"
	"gopkg.in/mgo.v2"
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
	err := c.Find(bson.M{"controlinfo.name": name, "controlinfo.tags": bson.M{"$all": tags}}).One(p)
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

func (db *DB) ListPackets(name string, tags []string) ([]*packet.Packet, error) {
	c := db.session.DB("jamesd").C("packets")
	query := bson.M{}
	if name != "" {
		query["controlinfo.name"] = name
	}
	if len(tags) > 0 {
		query["controlinfo.tags"] = bson.M{"$all": tags}
	}
	packets := make([]*packet.Packet, 0)
	err := c.Find(query).Select(bson.M{"controlinfo.name": 1, "controlinfo.tags": 1}).Sort("name", "tags").All(&packets)
	if err != nil {
		return nil, err
	}
	return packets, nil
}

func (db *DB) GetSatisfyingPackets(name string, tags []string) ([]*packet.Packet, error) {
	c := db.session.DB("jamesd").C("packets")
	query := bson.M{}
	if name != "" {
		query["controlinfo.name"] = name
	}
	mapFn := `
	function() {
		for (var i=0; i<this.controlinfo.tags.length; i++) {
			if (tags.indexOf(this.controlinfo.tags[i]) == -1) {
				return;
			}
		}
		emit( this._id, this );
	}
	`
	reduceFn := "function(key, values){ return values; }"
	job := &mgo.MapReduce{
		Map:    mapFn,
		Reduce: reduceFn,
		Scope:  bson.M{"tags": tags},
	}
	var result []struct {
		Id    string `bson:"_id"`
		Value *packet.Packet
	}
	_, err := c.Find(query).MapReduce(job, &result)
	if err != nil {
		return nil, err
	}
	packetList := make([]*packet.Packet, 0, len(result))
	for _, res := range result {
		packetList = append(packetList, res.Value)
	}
	return packetList, nil
}
