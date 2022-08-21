package geofencing

import (
	"context"
	"log"

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

func (gc *GeofencingCase) processMessageGeofencingDetect(ctx context.Context, msg repository.PubsubMessage) {
	log.Println(
		string(msg.GetMessage()),
	)
	msg.Ack()
}
