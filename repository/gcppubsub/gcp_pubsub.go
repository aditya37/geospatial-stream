package gpcppubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/aditya37/geospatial-stream/repository"
)

type gcpPubsub struct {
	client *pubsub.Client
}

func NewGCPPubsubRepo(client *pubsub.Client) (repository.PubsubManager, error) {
	return &gcpPubsub{
		client: client,
	}, nil
}

func (gp *gcpPubsub) Close() error {
	return gp.client.Close()
}
