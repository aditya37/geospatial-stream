package geofencing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	detect "github.com/aditya37/geofence-service/usecase"
	geofence_util "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/entity"
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
	internal_usecase "github.com/aditya37/geospatial-stream/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (gc *GeofencingCase) SubscribeGeofencingDetect(ctx context.Context, topicName, servicename string) error {
	if err := gc.pubsubManager.Subscribe(
		ctx,
		topicName,
		servicename,
		gc.processMessageGeofencingDetect,
	); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//
func (gc *GeofencingCase) processMessageGeofencingDetect(ctx context.Context, msg repository.PubsubMessage) {
	var payload detect.NotifyGeofencingPayload
	if err := json.Unmarshal(msg.GetMessage(), &payload); err != nil {
		geofence_util.Logger().Error(err)
		return
	}

	// set sseData...
	if isSetSSEData := gc.setSSEEventData(payload); !isSetSSEData {
		geofence_util.Logger().Info("channel nil")
	}

	// set to channel for stream websocket
	go gc.avgMobilityStream.SetToChannel(entity.AvgMobility{
		Inside: payload.Mobility.DailyAverage.Inside,
		Enter:  payload.Mobility.DailyAverage.Enter,
		Exit:   payload.Mobility.DailyAverage.Exit,
	})

	// writer for send/set Message Get Geofence By Channel
	gc.writeMessageGetGeofenceByChannel(payload)
	msg.Ack()
}

// writeMessageGetGeofenceByChannel...
func (gc *GeofencingCase) writeMessageGetGeofenceByChannel(payload detect.NotifyGeofencingPayload) {
	// get count device detected inside geofence area
	deviceCount, err := gc.trackingSvc.GetDeviceCounter(
		context.Background(),
		&emptypb.Empty{},
	)
	if err != nil {
		geofence_util.Logger().Error(err)
		return
	}

	// TODO: set to cache for speed optimize
	var detectedDevice []entity.DeviceType
	for _, d := range deviceCount.DetectDevice {
		detectedDevice = append(
			detectedDevice,
			entity.DeviceType{
				Type:     d.DeviceType,
				Count:    int(d.Count),
				DateTime: d.LastDetect,
			},
		)
	}

	msg := entity.PayloadMessageGetGeofenceByChann{
		Message:     fmt.Sprintf("Message From Channel: %s", payload.ChannelName),
		ChannelName: payload.ChannelName,
		Counter: entity.Counter{
			Enter:       payload.Mobility.DailyAverage.Enter,
			Inside:      payload.Mobility.DailyAverage.Inside,
			Exit:        payload.Mobility.DailyAverage.Exit,
			CountDevice: 0, // TODO: Get count device detected from geofence area
		},
		DeviceType: detectedDevice,
		Detect:     payload.Detect,
		Object:     payload.Object,
	}
	msgByte, _ := json.Marshal(msg)
	gc.chanStreamGeofenceByChannel.SendMessage(
		&channel.PayloadGetGeofenceByChann{
			Name: payload.ChannelName,
			Data: msgByte,
		},
	)
}

// setSSEEventData...
func (gc *GeofencingCase) setSSEEventData(payload detect.NotifyGeofencingPayload) bool {
	if theChan := gc.geofenceDetectChan.Get(payload.Type); theChan != nil {

		// assert response
		resp := internal_usecase.ResponseSSEDetectGeofence{
			Point:       payload.Object,
			Detect:      payload.Detect,
			ChannelName: payload.ChannelName,
			DeviceId:    payload.DeviceId,
		}
		j, _ := json.Marshal(resp)
		gc.geofenceDetectChan.Send(
			channel.Message{
				Type: payload.Type,
				Data: j,
			},
		)
		return true
	} else {
		return false
	}
}
