package service

import (
	"net/http"

	socket_deliv "github.com/aditya37/geospatial-stream/delivery/websocket"
	"github.com/gorilla/mux"
)

type httpHandler struct {
	muxRouter        *mux.Router
	geofencingSocket *socket_deliv.GeofencingWebsocketDeliver
}

func NewHttpHandler(
	geofencingSocket *socket_deliv.GeofencingWebsocketDeliver,
) *httpHandler {
	router := mux.NewRouter()
	return &httpHandler{
		muxRouter:        router,
		geofencingSocket: geofencingSocket,
	}
}

func (h *httpHandler) Handler() http.Handler {
	// health check
	h.muxRouter.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{'status':'up'}"))
	})

	// websocket...
	h.muxRouter.Methods(http.MethodGet).Path("/ws/geofencing/mobility/avg").HandlerFunc(h.geofencingSocket.StreamtMobilityAvg)
	return h.muxRouter
}
