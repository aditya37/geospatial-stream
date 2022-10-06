package geofencing

import (
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
)

type GeofencingCase struct {
	pubsubManager      repository.PubsubManager
	avgMobilityStream  *channel.StreamtMobilityAvg
	geofenceDetectChan *channel.ChannelStreamGeofencDetect
}

func NewGefencingUsecase(
	pubsubManager repository.PubsubManager,
	avgMobilityStream *channel.StreamtMobilityAvg,
	geofenceDetectChan *channel.ChannelStreamGeofencDetect,
) *GeofencingCase {
	return &GeofencingCase{
		pubsubManager:      pubsubManager,
		avgMobilityStream:  avgMobilityStream,
		geofenceDetectChan: geofenceDetectChan,
	}
}
