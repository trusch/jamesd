package db

import (
	"errors"

	"github.com/trusch/jamesd2/packet"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SavePacket saves a packet to db
func (db *DB) SavePacket(pack *packet.Packet) error {
	if _, err := pack.Hash(); err != nil {
		return err
	}
	if err := db.saveControlInfo(&pack.ControlInfo); err != nil {
		return err
	}
	if err := db.savePacketData(pack); err != nil {
		return err
	}
	return nil
}

// saveControlInfo saves controlinfo to db
func (db *DB) saveControlInfo(info *packet.ControlInfo) error {
	collection := db.db.C("controlinfo")
	_, err := collection.Upsert(bson.M{"name": info.Name, "labels": info.Labels}, info)
	return err
}

// savePacketData saves packets to db
func (db *DB) savePacketData(pack *packet.Packet) error {
	collection := db.db.C("packet")
	data, err := pack.ToData()
	if err != nil {
		return err
	}
	_, err = collection.Upsert(bson.M{"hash": pack.ControlInfo.Hash}, bson.M{"hash": pack.ControlInfo.Hash, "data": data})
	return err
}

// GetPacket gets a packet from db
func (db *DB) GetPacket(hash string) (*packet.Packet, error) {
	collection := db.db.C("packet")
	doc := &struct {
		Hash string
		Data []byte
	}{}
	err := collection.Find(bson.M{"hash": hash}).One(doc)
	if err != nil {
		return nil, err
	}
	pack, err := packet.NewFromData(doc.Data)
	if err != nil {
		return nil, err
	}
	_, err = pack.Hash()
	if err != nil {
		return nil, err
	}
	if hash != pack.ControlInfo.Hash {
		return nil, errors.New("packet hash mismatch")
	}
	return pack, nil
}

// DeletePacket deletes a packet
func (db *DB) DeletePacket(hash string) error {
	if err := db.db.C("packet").Remove(bson.M{"hash": hash}); err != nil {
		return err
	}
	if err := db.db.C("controlinfo").Remove(bson.M{"hash": hash}); err != nil {
		return err
	}
	return nil
}

// GetBestInfo returns controlinfo which doesnt contain a label which is not in the request
func (db *DB) GetBestInfo(name string, labels map[string]string) (*packet.ControlInfo, error) {
	collection := db.db.C("controlinfo")
	job := &mgo.MapReduce{
		Map: `function(){
      for(var key in this.labels){
        var reqVal = req[key];
        if(reqVal !== this.labels[key]){
          return;
        }
      }
      emit('', this)
    }`,
		Reduce: `function(key, values){
			var best = null;
			var maxLabelCount = 0;
			for(var idx in values){
				var labelCount = Object.keys(values[idx].labels).length;
				if(labelCount > maxLabelCount){
					best = values[idx];
					maxLabelCount = labelCount;
				}
			}
      return best;
    }`,
		Scope: map[string]interface{}{
			"req": labels,
		},
	}
	mapReduceResult := []struct{ Value *packet.ControlInfo }{}
	_, err := collection.Find(bson.M{"name": name}).MapReduce(job, &mapReduceResult)
	if err != nil {
		return nil, err
	}
	if len(mapReduceResult) == 0 {
		return nil, errors.New("no packet found")
	}
	return mapReduceResult[0].Value, nil
}

// GetInfos returns all controlinfos for packets for a given name
func (db *DB) GetInfos(name string) ([]*packet.ControlInfo, error) {
	collection := db.db.C("controlinfo")
	infos := make([]*packet.ControlInfo, 0, 8)
	if err := collection.Find(bson.M{"name": name}).All(&infos); err != nil {
		return nil, err
	}
	return infos, nil
}

// GetPacketNames returns a list of all distinct packet names
func (db *DB) GetPacketNames() ([]string, error) {
	collection := db.db.C("controlinfo")
	names := make([]string, 0, 32)
	if err := collection.Find(nil).Distinct("name", &names); err != nil {
		return nil, err
	}
	return names, nil
}
