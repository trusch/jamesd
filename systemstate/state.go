package systemstate

import (
	"encoding/gob"
	"os"
)

// SystemState contains info about the system, including version of installed apps
type SystemState struct {
	ID   string
	Arch string
	Apps []*AppInfo
}

// AppInfo contains name and version of the installed app
type AppInfo struct {
	Name    string
	Version string
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
	f, err := os.Open(path)
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

func (state *SystemState) SetAppVersion(name, version string) {
	for _, appInfo := range state.Apps {
		if appInfo.Name == name {
			appInfo.Version = version
			return
		}
	}
	state.Apps = append(state.Apps, &AppInfo{name, version})
}
