package geofencing

import (
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/aditya37/geospatial-tracking/proto"
)

type GeofencingCase struct {
	pubsubManager               repository.PubsubManager
	avgMobilityStream           *channel.StreamtMobilityAvg
	geofenceDetectChan          *channel.ChannelStreamGeofencDetect
	chanStreamGeofenceByChannel *channel.StreamGetGeofenceByChannel
	trackingSvc                 proto.GeotrackingClient
}

func NewGefencingUsecase(
	pubsubManager repository.PubsubManager,
	avgMobilityStream *channel.StreamtMobilityAvg,
	geofenceDetectChan *channel.ChannelStreamGeofencDetect,
	chanStreamGeofenceByChannel *channel.StreamGetGeofenceByChannel,
	trackingSvc proto.GeotrackingClient,
) *GeofencingCase {
	return &GeofencingCase{
		pubsubManager:               pubsubManager,
		avgMobilityStream:           avgMobilityStream,
		geofenceDetectChan:          geofenceDetectChan,
		chanStreamGeofenceByChannel: chanStreamGeofenceByChannel,
		trackingSvc:                 trackingSvc,
	}
}
