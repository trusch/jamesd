package state

import "github.com/trusch/jamesd/spec"

// State represents the state of a machine, i.e. which packets are installed (or should be installed)
type State struct {
	Apps []*App
}

// App represents a single installed packet
type App struct {
	*spec.App
	Hash string
}
