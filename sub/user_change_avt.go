package consumer

import (
	"context"
	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"log"
	"my-app/common"
	"my-app/common/asyncjob"
	"my-app/common/pubsub"
	"my-app/module/image"
)

type topicUserChangeAvt struct {
	s  sctx.ServiceContext
	ps pubsub.PubSub
}

func NewTopicUserChangeAvt(s sctx.ServiceContext, ps pubsub.PubSub) *topicUserChangeAvt {
	return &topicUserChangeAvt{s: s, ps: ps}
}

func (t *topicUserChangeAvt) setImgActAfterChangeAvt() error {
	ctx := context.Background()
	ch, _ := t.ps.Subscribe(ctx, common.TopicUserChangedAvt)

	dbCtx := t.s.MustGet(common.KeyGorm).(common.DbContext)

	for msg := range ch {
		mapData := msg.Data()
		imgId := uuid.MustParse(mapData["img_id"].(string))

		repo := image.NewRepo(dbCtx.GetDB())

		job := asyncjob.NewJob(func(ctx context.Context) error {
			return repo.SetImageStatusActivated(ctx, imgId)
		}, asyncjob.WithName("SetImgActAfterChangeAvt"))

		return asyncjob.NewGroup(false, job).Run(ctx)
	}

	return nil
}

//func (t *topicUserChangeAvt) doAnotherThing() error {
//	mapData := msg.Data()
//	imgId := uuid.MustParse(mapData["img_id"].(string))
//}

//func (t *topicUserChangeAvt) setImgActAfterChangeAvt2() error {
//	ctx := context.Background()
//	ch, _ := t.ps.Subscribe(ctx, common.TopicUserChangedAvt)
//
//	dbCtx := t.s.MustGet(common.KeyGorm).(common.DbContext)
//
//	for msg := range ch {
//		mapData := msg.Data()
//		imgId := uuid.MustParse(mapData["img_id"].(string))
//
//		repo := image.NewRepo(dbCtx.GetDB())
//
//		job := asyncjob.NewJob(func(ctx context.Context) error {
//			return repo.SetImageStatusActivated(ctx, imgId)
//		}, asyncjob.WithName("SetImgActAfterChangeAvt"))
//
//		return asyncjob.NewGroup(false, job).Run(ctx)
//	}
//
//	return nil
//}

func (t *topicUserChangeAvt) handlerSetImgActAfterChangeAvt(msg *pubsub.Message) error {
	mapData := msg.Data()
	imgId := uuid.MustParse(mapData["img_id"].(string))

	dbCtx := t.s.MustGet(common.KeyGorm).(common.DbContext)
	repo := image.NewRepo(dbCtx.GetDB())

	return repo.SetImageStatusActivated(context.Background(), imgId)
}

func (t *topicUserChangeAvt) doAnotherThing(msg *pubsub.Message) error {
	mapData := msg.Data()
	imgId := mapData["img_id"].(string)
	userId := mapData["user_id"].(string)

	log.Printf("user id: %s has changed avatar with image id: %s\n", userId, imgId)

	return nil
}

func (t *topicUserChangeAvt) Start() {
	ctx := context.Background()
	ch, _ := t.ps.Subscribe(ctx, common.TopicUserChangedAvt)

	for msg := range ch {
		job1 := asyncjob.NewJob(func(ctx context.Context) error {
			return t.handlerSetImgActAfterChangeAvt(msg)
		})

		job2 := asyncjob.NewJob(func(ctx context.Context) error {
			return t.doAnotherThing(msg)
		})

		asyncjob.NewGroup(true, job1, job2).Run(context.Background())
	}

}
