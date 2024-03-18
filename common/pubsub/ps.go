package pubsub

import (
	"context"
	"fmt"
	"time"
)

type PubSub interface {
	Publish(ctx context.Context, channel string, data *Message) error
	Subscribe(ctx context.Context, channel string) (ch <-chan *Message, close func())
}

type Message struct {
	id        string
	channel   string
	data      map[string]interface{}
	createdAt time.Time
}

func NewMessage(data map[string]interface{}) *Message {
	now := time.Now().UTC()

	return &Message{
		id:        fmt.Sprintf("%d", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

func (evt *Message) SetChannel(channel string) {
	evt.channel = channel
}

func (evt *Message) String() string {
	return fmt.Sprintf("Message %s value %v", evt.channel, evt.data)
}

func (evt *Message) Channel() string { return evt.channel }

func (evt *Message) Data() map[string]interface{} {
	return evt.data
}
