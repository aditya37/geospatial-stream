package channel

import (
	"log"
	"sync"
)

type (
	Client struct {
		Type    string
		Channel chan []byte
	}
	// contain message data
	Message struct {
		Type string
		Data []byte
	}
	ChannelStreamGeofencDetect struct {
		mutex sync.RWMutex
		// store clients...
		clients map[string]*Client
		// register client
		register chan *Client
		// send for store message send to client
		send chan Message
		// close
		close chan *Client
	}
)

func NewStreamGeofenceDetect() *ChannelStreamGeofencDetect {
	return &ChannelStreamGeofencDetect{
		mutex:    sync.RWMutex{},
		clients:  make(map[string]*Client),
		register: make(chan *Client),
		send:     make(chan Message),
		close:    make(chan *Client),
	}
}

// set data
func (cd *ChannelStreamGeofencDetect) Set(data *Client) {
	cd.mutex.Lock()
	cd.register <- data
	defer cd.mutex.Unlock()
}
func (cd *ChannelStreamGeofencDetect) Send(data Message) {
	cd.mutex.Lock()
	cd.send <- data
	defer cd.mutex.Unlock()
}
func (cd *ChannelStreamGeofencDetect) Get(data string) *Client {
	cd.mutex.RLock()
	d, ok := cd.clients[data]
	defer cd.mutex.RUnlock()
	if !ok {
		return nil
	}
	return d
}
func (cd *ChannelStreamGeofencDetect) Close(data *Client) {
	cd.mutex.Lock()
	cd.close <- data
	defer cd.mutex.Unlock()
}

// run...
func (cd *ChannelStreamGeofencDetect) Run() {
	for {
		select {
		// register client
		case cl := <-cd.register:
			cd.clients[cl.Type] = cl
			log.Printf("new client %d", len(cd.clients))
			// close and delete
		case cl := <-cd.close:
			if _, ok := cd.clients[cl.Type]; ok {
				delete(cd.clients, cl.Type)
			}
			log.Printf("%d", len(cd.clients))
		case msg := <-cd.send:
			if d, ok := cd.clients[msg.Type]; ok {
				select {
				case d.Channel <- msg.Data:
				default:
					close(d.Channel)
					delete(cd.clients, d.Type)
				}
			}
		}
	}
}
