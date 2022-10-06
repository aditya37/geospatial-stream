package geofencing

import (
	"context"
	"encoding/json"
	"log"

	detect "github.com/aditya37/geofence-service/usecase"
	geofence_util "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/entity"
	"github.com/aditya37/geospatial-stream/repository"
	"github.com/aditya37/geospatial-stream/repository/channel"
	internal_usecase "github.com/aditya37/geospatial-stream/usecase"
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

	msg.Ack()
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
