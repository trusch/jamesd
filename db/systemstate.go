package db

import (
	"github.com/trusch/jamesd/systemstate"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) SaveCurrentSystemState(state *systemstate.SystemState) error {
	return db.saveSystemState("systemstate", state)
}

func (db *DB) GetCurrentSystemState(name string) (*systemstate.SystemState, error) {
	return db.getSystemState("systemstate", name)
}

func (db *DB) saveSystemState(collection string, state *systemstate.SystemState) error {
	c := db.session.DB("jamesd").C(collection)
	_, err := c.Upsert(bson.M{"name": state.Name}, state)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getSystemState(collection, name string) (*systemstate.SystemState, error) {
	c := db.session.DB("jamesd").C(collection)
	state := &systemstate.SystemState{}
	err := c.Find(bson.M{"name": name}).One(state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (db *DB) ListSystems() ([]*systemstate.SystemState, error) {
	c := db.session.DB("jamesd").C("systemstate")
	systems := make([]*systemstate.SystemState, 0)
	result := &systemstate.SystemState{}
	iter := c.Find(nil).Select(bson.M{"name": 1, "timestamp": 1, "systemtags": 1}).Iter()
	for iter.Next(&result) {
		systems = append(systems, result)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return systems, nil
}
