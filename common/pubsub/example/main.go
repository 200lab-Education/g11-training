package main

import (
	"context"
	"fmt"
	"log"
	"my-app/common/pubsub"
	"time"
)

func main() {
	var broker pubsub.PubSub = pubsub.NewLocalPubSub("local-pubsub")

	const topic = "order.created"

	ch1, close1 := broker.Subscribe(context.Background(), topic)
	ch2, _ := broker.Subscribe(context.Background(), topic)

	go func() {
		for v := range ch1 {
			fmt.Println("Ch1:", v)
		}
		log.Println("Ch1 closed")
	}()

	go func() {
		for v := range ch2 {
			fmt.Println("Ch2:", v)
		}
	}()

	time.Sleep(time.Second * 3)
	broker.Publish(context.Background(), topic, pubsub.NewMessage(map[string]interface{}{"id": "ord-12389123"}))

	close1()

	time.Sleep(time.Second * 3)
	broker.Publish(context.Background(), topic, pubsub.NewMessage(map[string]interface{}{"id": "ord-12389124"}))

	time.Sleep(time.Second)

}
