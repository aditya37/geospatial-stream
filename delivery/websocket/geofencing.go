package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/aditya37/geospatial-stream/usecase/geofencing"
	gws "github.com/gorilla/websocket"
)

type GeofencingWebsocketDeliver struct {
	geofencingCase *geofencing.GeofencingCase
	upgrader       gws.Upgrader
	channelStream  *channel.StreamtMobilityAvg
}

func NewWebsocketGeofencing(
	geofencingCase *geofencing.GeofencingCase,
	channel *channel.StreamtMobilityAvg,
) *GeofencingWebsocketDeliver {
	socketUpgrader := gws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return &GeofencingWebsocketDeliver{
		geofencingCase: geofencingCase,
		upgrader:       socketUpgrader,
		channelStream:  channel,
	}
}

// StreamtMobilityAvg...
func (gwd *GeofencingWebsocketDeliver) StreamtMobilityAvg(w http.ResponseWriter, r *http.Request) {
	wsconn, err := gwd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ini error =>", err)
		return
	}

	go gwd.writeSocketMobilityAvg(wsconn)

}

// write or send to client socket...
func (gwd *GeofencingWebsocketDeliver) writeSocketMobilityAvg(conn *gws.Conn) {
	defer func() {
		log.Println("End Write...")
		conn.Close()
		gwd.channelStream.Close <- true
	}()
	for {
		select {
		case msg := <-gwd.channelStream.StreamAvgMobility:
			j, _ := json.Marshal(msg)
			if err := conn.WriteMessage(gws.TextMessage, j); err != nil {
				log.Println("err write", err)
				return
			}
		}
	}
}
