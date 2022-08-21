package repository

import (
	"context"
	"io"
)

type PubsubManager interface {
	io.Closer
	Subscribe(ctx context.Context, topicname, servicename string, callback interface{}) error
}

// PubsubMessage...
type PubsubMessage interface {
	GetMessage() []byte
	Ack()
}
