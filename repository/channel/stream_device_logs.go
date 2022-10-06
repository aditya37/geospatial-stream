package channel

import (
	"sync"

	"github.com/aditya37/geospatial-tracking/proto"
)

type ChannelStreamDeviceLogs struct {
	rwMutex     sync.RWMutex
	incomeData  chan []*proto.LogItem
	OutcomeData chan []*proto.LogItem
}

func NewStreamDeviceLogs(buffer int) *ChannelStreamDeviceLogs {
	if buffer == 0 {
		buffer = 1024
	}
	return &ChannelStreamDeviceLogs{
		incomeData:  make(chan []*proto.LogItem, buffer),
		rwMutex:     sync.RWMutex{},
		OutcomeData: make(chan []*proto.LogItem, buffer),
	}
}

func (dl *ChannelStreamDeviceLogs) SetToChannel(data []*proto.LogItem) {
	dl.rwMutex.Lock()
	defer dl.rwMutex.Unlock()
	dl.incomeData <- data
}

// run...
func (dl *ChannelStreamDeviceLogs) Run() {
	for {
		select {
		case data := <-dl.incomeData:
			dl.OutcomeData <- data
		}
	}
}
