package jamesd

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DB struct {
	session *mgo.Session
}

func NewDB(url string) (*DB, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &DB{session}, nil
}

type User struct {
	Name        string
	Permissions []string
}

func (db *DB) AddPermission(username, packageId string) error {
	c := db.session.DB("jamesd").C("users")
	var user User
	err := c.Find(bson.M{"name": username}).One(&user)
	if err != nil {
		user = User{Name: username, Permissions: []string{}}
	}
	for _, perm := range user.Permissions {
		if perm == packageId {
			return errors.New("user has already the permission")
		}
	}
	user.Permissions = append(user.Permissions, packageId)
	_, err = c.Upsert(bson.M{"name": username}, user)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RemovePermission(username, packageId string) error {
	c := db.session.DB("jamesd").C("users")
	var user User
	err := c.Find(bson.M{"name": username}).One(&user)
	if err != nil {
		user = User{Name: username, Permissions: []string{}}
	}
	newPermissions := make([]string, 0, len(user.Permissions))
	for _, perm := range user.Permissions {
		if perm != packageId {
			newPermissions = append(newPermissions, perm)
		}
	}
	user.Permissions = newPermissions
	_, err = c.Upsert(bson.M{"name": username}, user)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) CheckPermission(username, packageId string) bool {
	c := db.session.DB("jamesd").C("users")
	var user User
	err := c.Find(bson.M{"name": username}).One(&user)
	if err != nil {
		return false
	}
	for _, perm := range user.Permissions {
		if perm == packageId {
			return true
		}
	}
	return false
}

type Package struct {
	Name        string
	Arch        string
	Version     string
	Data        []byte
	Compression CompressionType
	Script      []byte
}

func (db *DB) AddPackage(pack *Package) error {
	c := db.session.DB("jamesd").C("packages")
	_, err := c.Upsert(pack, pack)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetPackage(name, arch, version string) (*Package, error) {
	pack := &Package{}
	c := db.session.DB("jamesd").C("packages")
	err := c.Find(bson.M{"name": name, "arch": arch, "version": version}).One(pack)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
