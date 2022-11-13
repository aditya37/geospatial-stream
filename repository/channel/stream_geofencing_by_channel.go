package channel

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type (
	// socket message
	PayloadGetGeofenceByChann struct {
		Name string
		Data []byte
	}
	// For store client
	GeofenceByChanClient struct {
		Name string
		// store instance of websocket
		WSConn *websocket.Conn
		// Pool...
		Pool *StreamGetGeofenceByChannel
	}

	// constructor....
	StreamGetGeofenceByChannel struct {
		mutex sync.RWMutex
		// store channel/client into map after register
		client map[string]*GeofenceByChanClient
		// register channel will store to map Client...
		register chan *GeofenceByChanClient
		// close
		// truncate or drop client from map
		close chan *GeofenceByChanClient
		// send message to specified client...
		send chan *PayloadGetGeofenceByChann
		// broadcast..
		// send broadcast message to all client
		broadcast chan *PayloadGetGeofenceByChann
	}
)

func NewStreamGeofenceByChannel() *StreamGetGeofenceByChannel {
	return &StreamGetGeofenceByChannel{
		mutex:     sync.RWMutex{},
		client:    make(map[string]*GeofenceByChanClient),
		register:  make(chan *GeofenceByChanClient),
		close:     make(chan *GeofenceByChanClient),
		send:      make(chan *PayloadGetGeofenceByChann),
		broadcast: make(chan *PayloadGetGeofenceByChann),
	}
}

// read Message from connection socket ..
func (cc *GeofenceByChanClient) ReadSocket() {
	defer func() {
		cc.WSConn.Close()
	}()

	for {
		if _, _, err := cc.WSConn.NextReader(); err != nil {
			log.Println(err)
			return
		}
	}

}

// register client...
func (gc *StreamGetGeofenceByChannel) Register(client *GeofenceByChanClient) {
	gc.mutex.Lock()
	gc.register <- client
	defer gc.mutex.Unlock()
}

// send message...
func (gc *StreamGetGeofenceByChannel) SendMessage(payload *PayloadGetGeofenceByChann) {
	gc.mutex.Lock()
	gc.send <- payload
	defer gc.mutex.Unlock()
}

func (gc *StreamGetGeofenceByChannel) SendBroadcast(payload *PayloadGetGeofenceByChann) {
	gc.mutex.Lock()
	gc.broadcast <- payload
	defer gc.mutex.Unlock()

}

// Run
// Run And stream....
func (gc *StreamGetGeofenceByChannel) Run() {
	for {
		select {
		case cl := <-gc.register:
			// store client from channel to map
			gc.client[cl.Name] = cl
		case client := <-gc.close:
			if _, ok := gc.client[client.Name]; ok {
				delete(gc.client, client.Name)
			}
		case bc := <-gc.broadcast:
			for _, c := range gc.client {
				c.WSConn.WriteMessage(
					1,
					bc.Data,
				)
			}
		case msg := <-gc.send:
			if client, ok := gc.client[msg.Name]; ok {
				client.WSConn.WriteMessage(
					1,
					msg.Data,
				)
			}
		}
	}
}
