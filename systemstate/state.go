package systemstate

import (
	"encoding/gob"
	"os"

	"github.com/trusch/jamesd/spec"
)

// SystemState contains info about the system, including version of installed apps
type SystemState struct {
	ID         string
	SystemTags []string
	Apps       []*spec.Entity
}

// AppInfo contains name and version of the installed app
type AppInfo struct {
	Name string
	Tags []string
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

func (state *SystemState) MarkAppInstalled(newApp *spec.Entity) {
	for _, info := range state.Apps {
		if info.Match(newApp) && newApp.Match(info) {
			// app already marked
			return
		}
	}
	state.Apps = append(state.Apps, newApp)
}

func (state *SystemState) MarkAppUninstalled(app *spec.Entity) {
	id := -1
	for idx, info := range state.Apps {
		if info.Match(app) && app.Match(info) {
			id = idx
			break
		}
	}
	if id != -1 {
		state.Apps = append(state.Apps[:id], state.Apps[id+1:]...)
	}
}
