package spec

// A Spec declares what should be installed where
// Target is a labels object specifying for which devices this spec matches
// Apps is a list of apps which should be installed, the also have a labels object
type Spec struct {
	ID     string
	Target map[string]string
	Apps   []*App
}

// App specifies an application to be installed.
// It consists of a Name and a labels object
type App struct {
	Name   string
	Labels map[string]string
}

// New returns a new spec
func New(id string) *Spec {
	return &Spec{ID: id, Target: make(map[string]string)}
}

// NewApp returns a new App
func NewApp(name string) *App {
	return &App{Name: name, Labels: make(map[string]string)}
}

// Clone creates a clone of a app
func (app *App) Clone() *App {
	res := NewApp(app.Name)
	for k, v := range app.Labels {
		res.Labels[k] = v
	}
	return res
}

// MergeLabels merges the given labels into the app spec
func (app *App) MergeLabels(labels map[string]string) {
	for k, v := range labels {
		app.Labels[k] = v
	}
}

// Clone creates a clone of a spec
func (s *Spec) Clone() *Spec {
	res := &Spec{Target: make(map[string]string)}
	for k, v := range s.Target {
		res.Target[k] = v
	}
	for _, app := range s.Apps {
		res.Apps = append(res.Apps, app.Clone())
	}
	return res
}
