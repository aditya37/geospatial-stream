package geofencing

import "github.com/aditya37/geospatial-stream/repository"

type GeofencingCase struct {
	pubsubManager repository.PubsubManager
}

func NewGefencingUsecase(
	pubsubManager repository.PubsubManager,
) *GeofencingCase {
	return &GeofencingCase{
		pubsubManager: pubsubManager,
	}
}
