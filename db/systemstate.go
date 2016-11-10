package db

import (
	"github.com/trusch/jamesd/systemstate"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) SaveCurrentSystemState(state *systemstate.SystemState) error {
	return db.saveSystemState("systemstate-current", state)
}

func (db *DB) SaveDesiredSystemState(state *systemstate.SystemState) error {
	return db.saveSystemState("systemstate-desired", state)
}

func (db *DB) GetCurrentSystemState(id string) (*systemstate.SystemState, error) {
	return db.getSystemState("systemstate-current", id)
}

func (db *DB) GetDesiredSystemState(id string) (*systemstate.SystemState, error) {
	return db.getSystemState("systemstate-desired", id)
}

func (db *DB) saveSystemState(collection string, state *systemstate.SystemState) error {
	c := db.session.DB("jamesd").C(collection)
	_, err := c.Upsert(bson.M{"id": state.ID}, state)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) getSystemState(collection, id string) (*systemstate.SystemState, error) {
	c := db.session.DB("jamesd").C(collection)
	state := &systemstate.SystemState{}
	err := c.Find(bson.M{"id": id}).One(state)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (db *DB) ListSystems() ([]string, error) {
	c := db.session.DB("jamesd").C("systemstate-current")
	systems := make([]string, 0)
	result := &systemstate.SystemState{}
	iter := c.Find(nil).Select(bson.M{"id": 1}).Iter()
	for iter.Next(&result) {
		systems = append(systems, result.ID)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return systems, nil
}
