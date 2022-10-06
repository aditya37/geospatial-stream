package device

import (
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
	"github.com/aditya37/geospatial-tracking/proto"
)

type DeviceManageUscase struct {
	pubsub           repository.PubsubManager
	trackingSvc      proto.GeotrackingClient
	streamDeviceLogs *channel.ChannelStreamDeviceLogs
}

func NewDeviceManagerUsecase(
	pubsub repository.PubsubManager,
	trackingSvc proto.GeotrackingClient,
	streamDeviceLogs *channel.ChannelStreamDeviceLogs,
) *DeviceManageUscase {
	return &DeviceManageUscase{
		pubsub:           pubsub,
		trackingSvc:      trackingSvc,
		streamDeviceLogs: streamDeviceLogs,
	}
}
