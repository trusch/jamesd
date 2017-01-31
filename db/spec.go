package db

import (
	"github.com/trusch/jamesd/spec"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SaveSpec saves a spec to db
func (db *DB) SaveSpec(spec *spec.Spec) error {
	collection := db.db.C("spec")
	_, err := collection.Upsert(bson.M{"id": spec.ID}, spec)
	return err
}

// GetSpec retrieves a spec from db
func (db *DB) GetSpec(id string) (*spec.Spec, error) {
	collection := db.db.C("spec")
	spec := &spec.Spec{}
	err := collection.Find(bson.M{"id": id}).One(&spec)
	return spec, err
}

// GetMergedSpec returns a merged specs of all matching specs in the db
func (db *DB) GetMergedSpec(labels map[string]string) (*spec.Spec, error) {
	collection := db.db.C("spec")
	job := &mgo.MapReduce{
		Map: `function(){
      for(var key in this.target) {
        var reqVal = req[key];
        if(reqVal !== this.target[key]){
          return;
        }
      }
      emit('', this)
    }`,
		Reduce: `function(key, specs){
      var res = {target: {}, apps: []};
      for(var idx in specs){
        var spec = specs[idx];
        for(var key in spec.target) {
          res.target[key] = spec.target[key];
        }
        for(var key in spec.apps) {
          res.apps.push(spec.apps[key]);
        }
      }
      return res;
    }`,
		Scope: map[string]interface{}{
			"req": labels,
		},
	}
	mapReduceResult := []struct{ Value *spec.Spec }{}
	_, err := collection.Find(nil).MapReduce(job, &mapReduceResult)
	if err != nil {
		return nil, err
	}
	if len(mapReduceResult) < 1 {
		return &spec.Spec{}, nil
	}
	return mapReduceResult[0].Value, nil
}

// DeleteSpec removes a spec from db
func (db *DB) DeleteSpec(id string) error {
	collection := db.db.C("spec")
	return collection.Remove(bson.M{"id": id})
}

// GetSpecs returns all specs
func (db *DB) GetSpecs() ([]*spec.Spec, error) {
	collection := db.db.C("spec")
	specs := []*spec.Spec{}
	err := collection.Find(nil).All(&specs)
	return specs, err
}
