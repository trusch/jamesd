package spec

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// A Spec specifies which packets should be installed on a set of targets
type Spec struct {
	Name   string `yaml:"name,omitempty"`
	Target *Entity
	Apps   []*Entity
}

func NewFromFile(filename string) (*Spec, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := &Spec{}
	return s, yaml.Unmarshal(bs, s)
}

// Merge merges the other spec into this one
func (spec *Spec) Merge(other *Spec) error {
	spec.Name = "merged-specs"
	if spec.Target.Name == "" {
		spec.Target.Name = other.Target.Name
	} else if other.Target.Name != "" && spec.Target.Name != other.Target.Name {
		return errors.New("cannot merge specs with different target names")
	}
	for _, otherTargetTag := range other.Target.Tags {
		if !spec.Target.CheckTag(otherTargetTag) {
			spec.Target.Tags = append(spec.Target.Tags, otherTargetTag)
		}
	}
	spec.Apps = append(spec.Apps, other.Apps...)
	return nil
}

// Match checks if this spec should be applied to a specific target
func (spec *Spec) Match(target *Entity) bool {
	return spec.Target.Match(target)
}

// An Entity describes either a target or a packet
type Entity struct {
	Name string `yaml:"name,omitempty"`
	Tags []string
}

// CheckTag returns true if the given tag is part of this entity
func (entity *Entity) CheckTag(tag string) bool {
	for _, t := range entity.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Match checks if another entity matches this one
// Therefore -> Name == other.Name && Tags <= other.Tags
func (entity *Entity) Match(other *Entity) bool {
	if entity.Name != other.Name {
		return false
	}
	for _, tag := range entity.Tags {
		if !other.CheckTag(tag) {
			return false
		}
	}
	return true
}
