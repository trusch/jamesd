package db

import (
	"errors"

	"github.com/trusch/jamesd/spec"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) UpsertSpec(specification *spec.Spec) error {
	c := db.session.DB("jamesd").C("specs")
	_, err := c.Upsert(bson.M{"name": specification.Name}, specification)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetSpec(name string) (*spec.Spec, error) {
	p := &spec.Spec{}
	c := db.session.DB("jamesd").C("specs")
	err := c.Find(bson.M{"name": name}).One(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (db *DB) RemoveSpec(name string) error {
	c := db.session.DB("jamesd").C("specs")
	return c.Remove(bson.M{"name": name})
}

func (db *DB) GetSpecs(targetName string, tags []string) ([]*spec.Spec, error) {
	c := db.session.DB("jamesd").C("specs")
	query := bson.M{}
	if targetName != "" {
		query["target.name"] = targetName
	}
	if len(tags) > 0 {
		query["target.tags"] = bson.M{"$all": tags}
	}
	specs := make([]*spec.Spec, 0)
	err := c.Find(query).Sort("target.name", "target.tags").All(&specs)
	if err != nil {
		return nil, err
	}
	return specs, nil
}

func (db *DB) GetSpecForTarget(name string, tags []string) (*spec.Spec, error) {
	c := db.session.DB("jamesd").C("specs")
	query := bson.M{
		"$or": []bson.M{
			bson.M{
				"target.tags": bson.M{
					"$not": bson.M{
						"$elemMatch": bson.M{
							"$nin": tags,
						},
					},
				},
			},
		},
	}
	if name != "" {
		query["$or"] = append(query["$or"].([]bson.M), bson.M{"target.name": name})
	}
	specs := []*spec.Spec{}
	err := c.Find(query).All(&specs)
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return &spec.Spec{}, errors.New("no matching spec found")
	}
	for i := 1; i < len(specs); i++ {
		specs[0].Merge(specs[i])
	}
	return specs[0], nil
}
