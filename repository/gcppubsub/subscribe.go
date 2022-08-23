package gpcppubsub

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/aditya37/geospatial-stream/repository"
	getenv "github.com/aditya37/get-env"
	"github.com/google/uuid"
)

// topic..
func (gp *gcpPubsub) createTopic(ctx context.Context, topic string) error {
	if _, err := gp.client.CreateTopic(ctx, topic); err != nil {
		return err
	}
	return nil
}

// get topic...
func (gp *gcpPubsub) getTopic(ctx context.Context, topicname string) (*pubsub.Topic, error) {
	topic := gp.client.Topic(topicname)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := gp.createTopic(ctx, topicname); err != nil {
			return nil, err
		}
		return nil, errors.New("topic not found")
	}
	return topic, nil
}

// createSubscription...
func (gp *gcpPubsub) createSubscription(ctx context.Context, servicename, topicname string) (*pubsub.Subscription, error) {
	topic, err := gp.getTopic(ctx, topicname)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err == nil {
		servicename = servicename + "." + id.String()
	}

	return gp.client.CreateSubscription(
		ctx,
		servicename,
		pubsub.SubscriptionConfig{
			Topic:               topic,
			RetainAckedMessages: getenv.GetBool("PUBSUB_RETAIN_ACKEDMSG", false),
			RetentionDuration: time.Duration(
				getenv.GetInt("PUBSUB_RETENTION_DURATION", 15) * int(time.Minute),
			),
		},
	)
}

// message callback...
type pubsubMessage struct {
	msg *pubsub.Message
}

func (pm pubsubMessage) GetMessage() []byte {
	return pm.msg.Data
}
func (pm pubsubMessage) Ack() {
	pm.msg.Ack()
}

func (gp *gcpPubsub) Subscribe(ctx context.Context, topicname, servicename string, callback interface{}) error {
	subs, err := gp.createSubscription(ctx, servicename, topicname)
	if err != nil {
		return err
	}
	if err := subs.Receive(
		ctx,
		func(ctx context.Context, m *pubsub.Message) {
			fn := callback.(func(context.Context, repository.PubsubMessage))
			fn(ctx, pubsubMessage{
				msg: m,
			})
		},
	); err != nil {
		return err
	}
	return nil
}
