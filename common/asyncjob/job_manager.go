package asyncjob

import (
	"context"
	"log"
	"my-app/common"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

//func (g *group) Run2(ctx context.Context) error {
//	errChan := make(chan error, len(g.jobs))
//
//	for i, _ := range g.jobs {
//		errChan <- g.runJob(ctx, g.jobs[i])
//	}
//
//	var err error
//
//	for i := 1; i <= len(g.jobs); i++ {
//		if v := <-errChan; v != nil {
//			err = v
//		}
//	}
//
//	return err
//}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for i := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				defer common.Recover()

				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(g.jobs[i])

			continue
		}

		job := g.jobs[i]

		err := g.runJob(ctx, job)

		if err != nil {
			return err
		}

		errChan <- err
		g.wg.Done()
	}

	g.wg.Wait()

	var err error

	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			return v
		}
	}

	return err
}

// Retry if needed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}
