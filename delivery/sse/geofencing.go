package sse

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aditya37/geospatial-stream/repository/channel"
)

type SSEGeofencingDeliver struct {
	stream  *channel.ChannelStreamGeofencDetect
	channel chan []byte
}

func NewSSEGeofencingDelivery(
	stream *channel.ChannelStreamGeofencDetect,
) *SSEGeofencingDeliver {
	channel := make(chan []byte, 1024)
	return &SSEGeofencingDeliver{
		stream:  stream,
		channel: channel,
	}
}

func (sd *SSEGeofencingDeliver) GeofenceDetectHandler(w http.ResponseWriter, r *http.Request) {
	key, ok := r.URL.Query()["type"]
	if !ok || len(key[0]) < 1 {
		return
	}

	// set keep alive
	sd.keepAliveHeader(w, r)

	// check channel
	if chl := sd.stream.Get(key[0]); chl == nil {
		sd.stream.Set(&channel.Client{
			Type:    key[0],
			Channel: sd.channel,
		})
	}

	// close
	notify := r.Context().Done()
	go func() {
		<-notify
		log.Printf("stop sse event : %s", key[0])
		sd.stream.Close(&channel.Client{
			Type: key[0],
		})
	}()

	flusher, _ := w.(http.Flusher)
	for {
		g := sd.stream.Get(key[0])
		if g != nil {
			fmt.Fprintf(w, "data: %s\n\n", <-sd.channel)
			// Flush the data immediatly instead of buffering it for later.
			flusher.Flush()
		}
	}
}

func (sd *SSEGeofencingDeliver) keepAliveHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}
