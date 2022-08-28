package channel

import (
	"sync"

	"github.com/aditya37/geospatial-stream/entity"
)

type StreamtMobilityAvg struct {
	rwMute            sync.RWMutex
	avgMobility       chan entity.AvgMobility
	StreamAvgMobility chan *entity.AvgMobility
	Close             chan bool
}

func NewStreamMobilityAvg() *StreamtMobilityAvg {
	return &StreamtMobilityAvg{
		rwMute:            sync.RWMutex{},
		avgMobility:       make(chan entity.AvgMobility),
		StreamAvgMobility: make(chan *entity.AvgMobility),
		Close:             make(chan bool),
	}
}

func (sm *StreamtMobilityAvg) SetToChannel(data entity.AvgMobility) {
	sm.rwMute.Lock()
	defer sm.rwMute.Unlock()
	sm.avgMobility <- data
}

func (sm *StreamtMobilityAvg) Run() {
	for {
		select {
		case data := <-sm.avgMobility:
			sm.StreamAvgMobility <- &data
		case c := <-sm.Close:
			if c == true {
				sm.StreamAvgMobility <- &entity.AvgMobility{}
			} else {
				continue
			}
		}
	}
}
