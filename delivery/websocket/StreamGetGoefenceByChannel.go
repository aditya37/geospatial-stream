package websocket

import (
	"log"
	"net/http"

	geofence_logger "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/gorilla/mux"
)

func (gwd *GeofencingWebsocketDeliver) StreamtGetGeofenceByChannel(w http.ResponseWriter, r *http.Request) {
	wsConn, err := gwd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		geofence_logger.Logger().Error(err)
		return
	}
	// parse param url....
	param := mux.Vars(r)
	channName, ok := param["channel_name"]
	if !ok {
		log.Println("channel_name param not set")
		return
	}

	cl := channel.GeofenceByChanClient{
		Name:   channName,
		WSConn: wsConn,
		Pool:   gwd.chanStreamGeofenceByChannel,
	}
	// register connection to channel
	gwd.chanStreamGeofenceByChannel.Register(&cl)

	// create goroutine for read message
	// for non blocking...
	go cl.ReadSocket()

}
