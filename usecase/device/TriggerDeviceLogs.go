package device

import (
	"context"
	"time"

	"github.com/aditya37/geospatial-tracking/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	geofence_logger "github.com/aditya37/geofence-service/util"
)

func (dm *DeviceManageUscase) TriggerDeviceLogs(ctx context.Context) ([]*proto.LogItem, error) {
	resp, err := dm.trackingSvc.GetDeviceLogs(
		ctx,
		&proto.RequestGetDeviceLogs{
			Paging: &proto.Paging{
				Page:        1,
				ItemPerPage: 5,
			},
			RecordedAt: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	)
	if err != nil {
		geofence_logger.Logger().Error(err)
		return nil, err
	}
	return resp.DeviceLogs, nil
}
