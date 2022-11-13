package service

import (
	"net/http"

	sse_deliv "github.com/aditya37/geospatial-stream/delivery/sse"
	socket_deliv "github.com/aditya37/geospatial-stream/delivery/websocket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type httpHandler struct {
	muxRouter        *mux.Router
	geofencingSocket *socket_deliv.GeofencingWebsocketDeliver
	deviceSocket     *socket_deliv.DeviceWebsocketDeliver
	geofencingSSE    *sse_deliv.SSEGeofencingDeliver
}

func NewHttpHandler(
	geofencingSocket *socket_deliv.GeofencingWebsocketDeliver,
	deviceSocket *socket_deliv.DeviceWebsocketDeliver,
	geofencingSSE *sse_deliv.SSEGeofencingDeliver,
) *httpHandler {
	router := mux.NewRouter()

	return &httpHandler{
		muxRouter:        router,
		geofencingSocket: geofencingSocket,
		deviceSocket:     deviceSocket,
		geofencingSSE:    geofencingSSE,
	}
}

func (h *httpHandler) Handler() http.Handler {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	// health check
	h.muxRouter.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{'status':'up'}"))
	})

	// websocket...
	h.muxRouter.Methods(http.MethodGet).Path("/ws/geofencing/mobility/avg").HandlerFunc(h.geofencingSocket.StreamtMobilityAvg)
	h.muxRouter.Methods(http.MethodGet).Path("/ws/geofencing/channel/{channel_name}").HandlerFunc(h.geofencingSocket.StreamtGetGeofenceByChannel)
	h.muxRouter.Methods(http.MethodGet).Path("/ws/devices/logs").HandlerFunc(h.deviceSocket.StreamDeviceLogs)

	// SSE Geofence service..
	geofencSrvcSse := h.muxRouter.PathPrefix("/sse/geofencing").Subrouter()
	geofencSrvcSse.Methods(http.MethodGet).Path("/").Queries("type", "").HandlerFunc(h.geofencingSSE.GeofenceDetectHandler)

	return corsHandler.Handler(h.muxRouter)
}
