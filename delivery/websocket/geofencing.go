package websocket

import (
	"encoding/json"
	"net/http"

	geofence_logger "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/aditya37/geospatial-stream/usecase/geofencing"
	gws "github.com/gorilla/websocket"
)

type GeofencingWebsocketDeliver struct {
	geofencingCase              *geofencing.GeofencingCase
	upgrader                    gws.Upgrader
	channelStream               *channel.StreamtMobilityAvg
	chanStreamGeofenceByChannel *channel.StreamGetGeofenceByChannel
}

func NewWebsocketGeofencing(
	geofencingCase *geofencing.GeofencingCase,
	channel *channel.StreamtMobilityAvg,
	chanStreamGeofenceByChannel *channel.StreamGetGeofenceByChannel,
) *GeofencingWebsocketDeliver {
	socketUpgrader := gws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return &GeofencingWebsocketDeliver{
		geofencingCase:              geofencingCase,
		upgrader:                    socketUpgrader,
		channelStream:               channel,
		chanStreamGeofenceByChannel: chanStreamGeofenceByChannel,
	}
}

// StreamtMobilityAvg...
func (gwd *GeofencingWebsocketDeliver) StreamtMobilityAvg(w http.ResponseWriter, r *http.Request) {
	wsconn, err := gwd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		geofence_logger.Logger().Error(err)
		return
	}
	// Reader...
	go gwd.readSocketMobilityAvg(wsconn)
	// writer...
	go gwd.writeSocketMobilityAvg(wsconn)

}

// read..
func (gwd *GeofencingWebsocketDeliver) readSocketMobilityAvg(conn *gws.Conn) {
	defer func() {
		// clear resource...
		conn.Close()
		gwd.channelStream.Close <- true
	}()
	for {
		if _, _, err := conn.NextReader(); err != nil {
			geofence_logger.Logger().Error(err)
			return
		}
	}
}

// write or send to client socket...
func (gwd *GeofencingWebsocketDeliver) writeSocketMobilityAvg(conn *gws.Conn) {
	// start write
	gwd.channelStream.Close <- false

	defer func() {
		geofence_logger.Logger().Info("CLose write...")
		conn.Close()
	}()

	for {
		select {
		case msg := <-gwd.channelStream.StreamAvgMobility:
			j, _ := json.Marshal(msg)
			if err := conn.WriteMessage(gws.TextMessage, j); err != nil {
				geofence_logger.Logger().Error(err)
				return
			}

		}
	}
}
