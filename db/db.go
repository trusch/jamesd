package db

import mgo "gopkg.in/mgo.v2"

// DB wraps a mongodb connection
type DB struct {
	session *mgo.Session
	db      *mgo.Database
}

// New creates a new db object
func New(uri string) (*DB, error) {
	info, err := mgo.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	return &DB{
		session: session,
		db:      session.DB(info.Database),
	}, nil
}

// Drop drops the entire database
func (db *DB) Drop() error {
	return db.db.DropDatabase()
}

// Close closes the db
func (db *DB) Close() error {
	db.session.Close()
	return nil
}
