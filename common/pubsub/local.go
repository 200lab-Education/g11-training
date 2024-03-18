package pubsub

import (
	"context"
	"log"
	"my-app/common"
	"sync"
)

type localPubSub struct {
	name         string
	messageQueue chan *Message
	mapChannel   map[string][]chan *Message
	locker       *sync.RWMutex
}

func NewLocalPubSub(name string) *localPubSub {
	pb := &localPubSub{
		name:         name,
		messageQueue: make(chan *Message, 10000),
		mapChannel:   make(map[string][]chan *Message),
		locker:       new(sync.RWMutex),
	}

	pb.run()

	return pb
}

func (ps *localPubSub) Publish(ctx context.Context, topic string, data *Message) error {
	data.SetChannel(topic)

	go func() {
		defer common.Recover()
		ps.messageQueue <- data
		log.Println("New message published:", data.String())
	}()

	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, topic string) (ch <-chan *Message, unsubscribe func()) {
	c := make(chan *Message)

	ps.locker.Lock()

	if val, ok := ps.mapChannel[topic]; ok {
		val = append(ps.mapChannel[topic], c)
		ps.mapChannel[topic] = val
	} else {
		ps.mapChannel[topic] = []chan *Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[topic]; ok {
			for i := range chans {
				if chans[i] == c {
					// remove element at index in chans
					// [1,2,3,4,5] //  i = 3
					// [1,2,3] (arr[:i])
					// [5] (arr[i+1:])
					// [1,2,3,5]
					chans = append(chans[:i], chans[i+1:]...)

					ps.locker.Lock()
					ps.mapChannel[topic] = chans
					ps.locker.Unlock()

					close(c)
					break
				}
			}
		}
	}
}

func (ps *localPubSub) run() error {
	go func() {
		defer common.Recover()
		for {
			mess := <-ps.messageQueue
			log.Println("Message dequeue:", mess.String())

			if subs, ok := ps.mapChannel[mess.Channel()]; ok {
				for i := range subs {
					go func(c chan *Message) {
						defer common.Recover()
						c <- mess
						//f(mess)
					}(subs[i])
				}
			}
			//else {
			//	ps.messageQueue <- mess
			//}
		}
	}()

	return nil
}

func (ps *localPubSub) GetPrefix() string {
	return ps.name
}

func (ps *localPubSub) Get() interface{} {
	return ps
}

func (ps *localPubSub) Name() string {
	return ps.name
}

func (ps *localPubSub) InitFlags() {
}

func (ps *localPubSub) Configure() error {
	return nil
}

func (ps *localPubSub) Run() error {
	return nil
}

func (ps *localPubSub) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}
