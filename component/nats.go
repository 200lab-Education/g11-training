package component

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/nats-io/nats.go"
	sctx "github.com/viettranx/service-context"
	"log"
	"my-app/common/pubsub"
	"time"
)

type psNATS struct {
	id         string
	url        string
	connection *nats.Conn
}

func NewNATSComponent(id string) *psNATS {
	return &psNATS{
		id: id,
	}
}

func (n *psNATS) ID() string {
	return n.id
}

func (n *psNATS) InitFlags() {
	flag.StringVar(&n.url, "nats-url", nats.DefaultURL, "URL of NATS service")
}

func (n *psNATS) Activate(context sctx.ServiceContext) error {
	log.Println("Connecting to NATS service...")
	conn, err := nats.Connect(n.url, n.setupConnOptions([]nats.Option{})...)

	if err != nil {
		log.Println(err)
	}

	log.Println("Connected to NATS service.")

	n.connection = conn

	return nil
}

func (n *psNATS) Stop() error {
	n.connection.Close()
	return nil
}

func (n *psNATS) setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Println("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Println("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Println("Exiting: %v", nc.LastError())
	}))

	return opts
}

func (n *psNATS) Publish(ctx context.Context, channel string, data *pubsub.Message) error {
	msgData, err := json.Marshal(data.Data())

	if err != nil {
		log.Println(err)
		return err
	}

	if err := n.connection.Publish(channel, msgData); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (n *psNATS) Subscribe(ctx context.Context, channel string) (ch <-chan *pubsub.Message, close func()) {
	msgChan := make(chan *pubsub.Message)

	//go func() {}()
	sub, err := n.connection.Subscribe(channel, func(msg *nats.Msg) {
		msgData := make(map[string]interface{})

		_ = json.Unmarshal(msg.Data, &msgData)

		appMsg := pubsub.NewMessage(msgData)
		appMsg.SetChannel(channel)

		msgChan <- appMsg

	})

	if err != nil {
		log.Println(err)
	}

	return msgChan, func() {
		_ = sub.Unsubscribe()
	}
}
