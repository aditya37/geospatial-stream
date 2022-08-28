package device

import (
	"context"
	"encoding/json"

	"github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-tracking/proto"
)

func (dm *DeviceManageUscase) SubscribeDeviceLogs(ctx context.Context, topicName, servicename string) error {
	if err := dm.pubsub.Subscribe(
		ctx,
		topicName,
		servicename,
		dm.callbackStreamDeviceLogs,
	); err != nil {
		util.Logger().Error(err)
		return err
	}
	return nil
}

func (dm *DeviceManageUscase) callbackStreamDeviceLogs(ctx context.Context, msg repository.PubsubMessage) {
	var payload []*proto.LogItem
	if err := json.Unmarshal(msg.GetMessage(), &payload); err != nil {
		util.Logger().Error(err)
		return
	}

	dm.streamDeviceLogs.SetToChannel(payload)
	msg.Ack()
}
