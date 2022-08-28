package geofencing

import (
	"context"
	"encoding/json"
	"log"

	detect "github.com/aditya37/geofence-service/usecase"
	geofence_util "github.com/aditya37/geofence-service/util"
	"github.com/aditya37/geospatial-stream/entity"
	"github.com/aditya37/geospatial-stream/repository"
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

	// set to channel for stream websocket
	go gc.avgMobilityStream.SetToChannel(entity.AvgMobility{
		Inside: payload.Mobility.DailyAverage.Inside,
		Enter:  payload.Mobility.DailyAverage.Enter,
		Exit:   payload.Mobility.DailyAverage.Exit,
	})

	msg.Ack()
}
