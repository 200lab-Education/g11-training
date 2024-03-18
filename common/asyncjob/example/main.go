package main

import (
	"context"
	"log"
	"my-app/common/asyncjob"
	"time"
)

func main() {
	j1 := asyncjob.NewJob(func(ctx context.Context) error {
		log.Println("doing something at J1")
		return nil
	}, asyncjob.WithName("J1"))

	j2 := asyncjob.NewJob(func(ctx context.Context) error {
		log.Println("doing something at J2")
		return nil
	}, asyncjob.WithName("J2"), asyncjob.WithRetriesDuration([]time.Duration{time.Second * 5}))

	jm := asyncjob.NewGroup(true, j1, j2)

	if err := jm.Run(context.Background()); err != nil {
		log.Println(err)
	}

	//if err := j1.Execute(context.Background()); err != nil {
	//	log.Println(err)
	//
	//	for {
	//		err := j1.Retry(context.Background())
	//
	//		if err == nil || j1.State() == asyncjob.StateRetryFailed {
	//			break
	//		}
	//	}
	//}
}
