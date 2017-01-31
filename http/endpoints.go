package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trusch/jamesd2/packet"
	"github.com/trusch/jamesd2/spec"
	"github.com/trusch/jamesd2/state"
)

func (srv *server) listPackets(w http.ResponseWriter, r *http.Request) {
	packetNames, err := srv.db.GetPacketNames()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := make(map[string][]*packet.ControlInfo)
	for _, name := range packetNames {
		infos, err := srv.db.GetInfos(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res[name] = infos
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(res)
}

func (srv *server) postPacket(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pack, err := packet.NewFromData(bs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = srv.db.SavePacket(pack)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(pack.ControlInfo)
}

func (srv *server) deletePacket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	err := srv.db.DeletePacket(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (srv *server) getPacketData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	pack, err := srv.db.GetPacket(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	bs, err := pack.ToData()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bs)
}

func (srv *server) getPacketInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	pack, err := srv.db.GetPacket(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(pack.ControlInfo)
}

func (srv *server) computePacketList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	labels := make(map[string]string)
	err := decoder.Decode(labels)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err := srv.db.GetMergedSpec(labels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	desiredState := &state.State{}
	for _, app := range s.Apps {
		app.MergeLabels(labels)
		info, err := srv.db.GetBestInfo(app.Name, app.Labels)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		desired := &state.App{
			App: &spec.App{
				Name:   info.Name,
				Labels: info.Labels,
			},
			Hash: info.Hash,
		}
		desiredState.Apps = append(desiredState.Apps, desired)
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(desiredState)
}

func (srv *server) listSpecs(w http.ResponseWriter, r *http.Request) {
	specs, err := srv.db.GetSpecs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(specs)
}

func (srv *server) postSpec(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	s := &spec.Spec{}
	err := decoder.Decode(s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = srv.db.SaveSpec(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(s)
}

func (srv *server) computeSpec(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	labels := make(map[string]string)
	err := decoder.Decode(labels)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s, err := srv.db.GetMergedSpec(labels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(s)
}

func (srv *server) putSpec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	decoder := json.NewDecoder(r.Body)
	clientSpec := &spec.Spec{}
	err := decoder.Decode(clientSpec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = srv.db.GetSpec(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = srv.db.SaveSpec(clientSpec)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(clientSpec)
}

func (srv *server) deleteSpec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := srv.db.DeleteSpec(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (srv *server) getSpec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	s, err := srv.db.GetSpec(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	encoder.Encode(s)
}
