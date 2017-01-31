package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/trusch/jamesd/db"
)

type server struct {
	handler *mux.Router
	db      *db.DB
}

func (srv *server) buildEndpoint() {
	router := mux.NewRouter()

	packetRouter := router.PathPrefix("/packet").Subrouter().StrictSlash(true)
	packetRouter.Path("/").Methods("GET").HandlerFunc(srv.listPackets)
	packetRouter.Path("/").Methods("POST").HandlerFunc(srv.postPacket)
	packetRouter.Path("/compute").Methods("POST").HandlerFunc(srv.computePacketList)
	packetRouter.Path("/{hash}").Methods("DELETE").HandlerFunc(srv.deletePacket)
	packetRouter.Path("/{hash}/data").Methods("GET").HandlerFunc(srv.getPacketData)
	packetRouter.Path("/{hash}/info").Methods("GET").HandlerFunc(srv.getPacketInfo)

	specRouter := router.PathPrefix("/spec").Subrouter().StrictSlash(true)
	specRouter.Path("/").Methods("GET").HandlerFunc(srv.listSpecs)
	specRouter.Path("/").Methods("POST").HandlerFunc(srv.postSpec)
	specRouter.Path("/compute").Methods("POST").HandlerFunc(srv.computeSpec)
	specRouter.Path("/{id}").Methods("GET").HandlerFunc(srv.getSpec)
	specRouter.Path("/{id}").Methods("PUT").HandlerFunc(srv.putSpec)
	specRouter.Path("/{id}").Methods("DELETE").HandlerFunc(srv.deleteSpec)

	srv.handler = router
}

func ListenAndServe(db *db.DB, addr string) error {
	server := &server{db: db}
	server.buildEndpoint()
	return http.ListenAndServe(addr, server.handler)
}
