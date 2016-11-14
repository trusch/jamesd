package systemstate

import (
	"encoding/gob"
	"os"
)

// SystemState contains info about the system, including version of installed apps
type SystemState struct {
	ID         string
	SystemTags []string
	Apps       []*AppInfo
}

// AppInfo contains name and version of the installed app
type AppInfo struct {
	Name string
	Tags []string
}

func (info *AppInfo) Equals(other *AppInfo) bool {
	if info.Name != other.Name || len(info.Tags) != len(other.Tags) {
		return false
	}
	ownTags := make(map[string]bool)
	for _, tag := range info.Tags {
		ownTags[tag] = true
	}
	for _, tag := range other.Tags {
		if _, ok := ownTags[tag]; !ok {
			return false
		}
	}
	return true
}

func (info *AppInfo) IsSuperSetOf(other *AppInfo) bool {
	if info.Name != other.Name {
		return false
	}
	for _, tag := range info.Tags {
		foundTag := false
		for _, otherTag := range other.Tags {
			if tag == otherTag {
				foundTag = true
				break
			}
		}
		if !foundTag {
			return false
		}
	}
	return true
}

func (state *SystemState) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(f)
	return decoder.Decode(state)
}

func (state *SystemState) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(f)
	return encoder.Encode(state)
}

func NewFromFile(path string) (*SystemState, error) {
	state := &SystemState{}
	return state, state.Load(path)
}

func (state *SystemState) MarkAppInstalled(newApp *AppInfo) {
	for _, info := range state.Apps {
		if info.Equals(newApp) {
			// app already marked
			return
		}
	}
	state.Apps = append(state.Apps, newApp)
}

func (state *SystemState) MarkAppUninstalled(app *AppInfo) {
	id := -1
	for idx, info := range state.Apps {
		if info.Equals(app) {
			id = idx
			break
		}
	}
	if id != -1 {
		state.Apps = append(state.Apps[:id], state.Apps[id+1:]...)
	}
}
