package geofencing

import (
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
)

type GeofencingCase struct {
	pubsubManager     repository.PubsubManager
	avgMobilityStream *channel.StreamtMobilityAvg
}

func NewGefencingUsecase(
	pubsubManager repository.PubsubManager,
	avgMobilityStream *channel.StreamtMobilityAvg,
) *GeofencingCase {
	return &GeofencingCase{
		pubsubManager:     pubsubManager,
		avgMobilityStream: avgMobilityStream,
	}
}
