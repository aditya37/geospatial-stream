package infra

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

var (
	pubsubInstance  *pubsub.Client = nil
	pubsubSingleton sync.Once
	returnErr       error
)

func NewGCPPubsub(ctx context.Context, projectid string, opts ...option.ClientOption) error {
	pubsubSingleton.Do(func() {
		client, err := pubsub.NewClient(
			ctx,
			projectid,
			opts...,
		)
		if err != nil {
			returnErr = err
			return
		}
		pubsubInstance = client
	})
	if returnErr != nil {
		return returnErr
	}
	return nil
}
func GetGCPPubsubInstance() *pubsub.Client {
	return pubsubInstance
}
